package ghaprofiler

import (
	"context"
	"log"
	"os"
	"regexp"
	"sort"

	"github.com/google/go-github/v32/github"
	"github.com/jessevdk/go-flags"
	"golang.org/x/sync/errgroup"
)

func StartCLI(ctx context.Context, args []string) {
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

	if config.Verbose {
		log.Printf("config=%v", configTomlPath)
		log.Print(config.Dump())
	}

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

	if config.Verbose {
		log.Println("ListWorkflowRunsByFileName start")
	}
	workflowRuns, _, err := client.ListWorkflowRunsByFileName(ctx, config.Owner, config.Repository, config.WorkflowFileName, listWorkflowRunsOpts)
	if err != nil {
		log.Fatal(err)
	}
	if config.Verbose {
		log.Println("ListWorkflowRunsByFileName finish")
	}

	jobsByJobName := NewJobsByJobNameMap()
	eg := new(errgroup.Group)
	sem := make(chan struct{}, config.Concurrency)

	for _, run := range workflowRuns.WorkflowRuns {
		sem <- struct{}{}
		eg.Go(func() error {
			defer func() {
				<-sem
			}()
			if config.Verbose {
				log.Printf("ListWorkflowJobs start: run_id=%d", *run.ID)
			}
			jobs, _, err := client.ListWorkflowJobs(ctx, config.Owner, config.Repository, *run.ID, nil)
			if err != nil {
				return err
			}
			if config.Verbose {
				log.Printf("ListWorkflowJobs finish: run_id=%d", *run.ID)
			}

			for _, job := range jobs.Jobs {
				jobName := *job.Name
				if !jobNameRegex.MatchString(jobName) {
					continue
				}
				if config.Verbose {
					log.Printf("Job name (before replacement): %#v", jobName)
				}
				for _, rule := range config.Replace {
					jobName = rule.Apply(jobName)
				}
				if config.Verbose {
					log.Printf("Job name (after replacement): %#v", jobName)
				}
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
