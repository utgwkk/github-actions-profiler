package ghaprofiler

import (
	"context"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v32/github"
	"github.com/jessevdk/go-flags"
	"golang.org/x/sync/errgroup"
)

type CLI struct {
	verbose bool
}

func NewCLI() *CLI {
	return &CLI{}
}

func (cli *CLI) SetVerbosity(verbose bool) {
	cli.verbose = verbose
}

func (cli *CLI) logVerbose(s interface{}) {
	if !cli.verbose {
		return
	}
	log.Print(s)
}

func (cli *CLI) loglnVerbose(s interface{}) {
	if !cli.verbose {
		return
	}
	log.Println(s)
}

func (cli *CLI) logfVerbose(format string, args ...interface{}) {
	if !cli.verbose {
		return
	}
	log.Printf(format, args...)
}

func (cli *CLI) overrideRepositoryFromCWD(config *ProfileConfig) {
	if config.Owner != "" && config.Repository != "" {
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	repo, err := git.PlainOpen(cwd)
	if err != nil {
		cli.loglnVerbose(err)
		return
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		cli.loglnVerbose(err)
		return
	}
	urls := remote.Config().URLs
	if len(urls) == 0 {
		return
	}
	url := urls[0]
	splitted := strings.Split(url, "/")
	if len(splitted) <= 1 {
		return
	}
	owner, repoName := splitted[len(splitted)-2], splitted[len(splitted)-1]
	repoName = strings.Replace(repoName, ".git", "", 1)
	if lastColonPos := strings.LastIndex(owner, ":"); lastColonPos != -1 {
		// git@github.com:owner/repo.git
		owner = owner[lastColonPos:]
	}
	config.Owner = owner
	config.Repository = repoName
}

func (cli *CLI) Start(ctx context.Context, args []string) {
	var config *ProfileConfig = DefaultProfileConfig()
	var configFromArgs ProfileConfigCLIArgs
	args, err := flags.ParseArgs(&configFromArgs, args)
	if err != nil {
		// flags.ParseArgs() outputs error message, so discarding it here...
		return
	}

	var configTomlPath string
	if configFromArgs.ConfigPath != nil {
		configTomlPath = *configFromArgs.ConfigPath
		configFromTOML, err := LoadConfigFromTOML(configTomlPath)
		if err != nil {
			log.Fatalf("Failed to load %s: %v", configTomlPath, err)
		}
		config = OverrideCLIArgs(configFromTOML, &configFromArgs)
	} else {
		config = OverrideCLIArgs(DefaultProfileConfig(), &configFromArgs)
	}
	cli.overrideRepositoryFromCWD(config)
	cli.SetVerbosity(config.Verbose)

	cli.logfVerbose("config=%v", configTomlPath)
	cli.logVerbose(config.Dump())

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	jobNameRegex, err := regexp.Compile(config.JobNameRegexp)
	if err != nil {
		log.Fatal(err)
	}

	client := NewClientWithConfig(ctx, &ClientConfig{
		AccessToken:    config.AccessToken,
		Cache:          config.Cache,
		CacheDirectory: config.CacheDirectory,
	})

	listWorkflowRunsOpts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			PerPage: config.NumberOfJob,
		},
	}

	cli.loglnVerbose("ListWorkflowRunsByFileName start")
	workflowRuns, _, err := client.ListWorkflowRunsByFileName(ctx, config.Owner, config.Repository, config.WorkflowFileName, listWorkflowRunsOpts)
	if err != nil {
		log.Fatal(err)
	}
	cli.loglnVerbose("ListWorkflowRunsByFileName finish")

	jobsByJobName := NewJobsByJobNameMap()
	eg := new(errgroup.Group)
	sem := make(chan struct{}, config.Concurrency)

	for _, run := range workflowRuns.WorkflowRuns {
		sem <- struct{}{}
		eg.Go(func() error {
			defer func() {
				<-sem
			}()
			cli.logfVerbose("ListWorkflowJobs start: run_id=%d", *run.ID)
			jobs, _, err := client.ListWorkflowJobs(ctx, config.Owner, config.Repository, *run.ID, nil)
			if err != nil {
				return err
			}
			cli.logfVerbose("ListWorkflowJobs finish: run_id=%d", *run.ID)

			for _, job := range jobs.Jobs {
				jobName := *job.Name
				if !jobNameRegex.MatchString(jobName) {
					continue
				}
				cli.logfVerbose("Job name (before replacement): %#v", jobName)
				for _, rule := range config.Replace {
					jobName = rule.Apply(jobName)
				}
				cli.logfVerbose("Job name (after replacement): %#v", jobName)
				jobsByJobName.Append(jobName, job)
			}
			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		log.Fatal(err)
	}

	profileResult := make(map[string][]*TaskStepProfile)

	for jobName, jobs := range jobsByJobName.Iterate() {
		if len(jobs) == 0 {
			continue
		}

		var steps []*github.TaskStep
		for _, job := range jobs {
			steps = append(steps, job.Steps...)
		}

		stepProfile, err := ProfileTaskStep(steps)
		if err != nil {
			log.Fatal(err)
		}
		err = SortProfileBy(stepProfile, config.SortBy)
		if err != nil {
			log.Fatal(err)
		}

		// reverse slice
		if config.Reverse {
			reversedStepProfile := make(TaskStepProfileResult, len(stepProfile))
			for i := 0; i < len(stepProfile); i++ {
				j := len(stepProfile) - i - 1
				reversedStepProfile[i] = stepProfile[j]
			}
			profileResult[jobName] = reversedStepProfile
		} else {
			profileResult[jobName] = stepProfile
		}
	}

	var formatterInputJobNames []string
	for k := range profileResult {
		formatterInputJobNames = append(formatterInputJobNames, k)
	}
	sort.Strings(formatterInputJobNames)

	var profileFormatterInput ProfileInput
	for _, jobName := range formatterInputJobNames {
		result := profileResult[jobName]
		profileFormatterInput = append(profileFormatterInput, &ProfileForFormatter{
			Name:    jobName,
			Profile: result,
		})
	}

	WriteWithFormat(os.Stdout, profileFormatterInput, config.Format)
}
