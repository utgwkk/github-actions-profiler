package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"sort"

	"github.com/google/go-github/v32/github"
	"github.com/jessevdk/go-flags"
	ghaprofiler "github.com/utgwkk/github-actions-profiler"
	"golang.org/x/sync/errgroup"
)

var config *ghaprofiler.ProfileConfig = ghaprofiler.DefaultProfileConfig()

const accessTokenEnvVariableName = "GITHUB_ACTIONS_PROFILER_TOKEN"

func main() {
	ctx := context.Background()

	var configFromArgs ghaprofiler.ProfileConfigCLIArgs
	args := os.Args[1:]
	args, err := flags.ParseArgs(&configFromArgs, args)
	if err != nil {
		// flags.ParseArgs() outputs error message, so discarding it here...
		return
	}

	var configTomlPath string
	if configFromArgs.ConfigPath != nil {
		configTomlPath = *configFromArgs.ConfigPath
		configFromTOML, err := ghaprofiler.LoadConfigFromTOML(configTomlPath)
		if err != nil {
			log.Fatalf("Failed to load %s: %v", configTomlPath, err)
		}
		// TODO: Override config with CLI arguments when they are given
		config = ghaprofiler.OverrideCLIArgs(configFromTOML, &configFromArgs)
	} else {
		config = ghaprofiler.OverrideCLIArgs(ghaprofiler.DefaultProfileConfig(), &configFromArgs)
	}

	if config.Verbose {
		log.Printf("config=%v\n", configTomlPath)
		log.Printf("count=%v\n", config.Count)
		log.Printf("format=%v\n", config.Format)
		log.Printf("job-name-regexp=%v\n", config.JobNameRegexp)
		log.Printf("owner=%v\n", config.Owner)
		log.Printf("repo=%v\n", config.Repository)
		log.Printf("reverse=%v\n", config.Reverse)
		log.Printf("sort=%v\n", config.SortBy)
		// We don't write out token
		log.Printf("workflow_file=%v\n", config.WorkflowFileName)
		log.Printf("replace=%#v", config.Replace)
		log.Printf("cache=%v", config.Cache)
		log.Printf("cache_directory=%v", config.CacheDirectory)
	}

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	if config.AccessToken == "" {
		accessTokenFromEnv := os.Getenv(accessTokenEnvVariableName)
		if accessTokenFromEnv != "" {
			config.AccessToken = accessTokenFromEnv
		}
	}

	jobNameRegex, err := regexp.Compile(config.JobNameRegexp)
	if err != nil {
		log.Fatal(err)
	}

	client := ghaprofiler.NewClientWithConfig(ctx, &ghaprofiler.ClientConfig{
		AccessToken:    config.AccessToken,
		Cache:          config.Cache,
		CacheDirectory: config.CacheDirectory,
	})

	listWorkflowRunsOpts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			PerPage: config.Count,
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
	// TODO: make configurable
	sem := make(chan struct{}, 4)

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

	profileResult := make(map[string][]*ghaprofiler.TaskStepProfile)

	for jobName, jobs := range jobsByJobName.Iterate() {
		if len(jobs) == 0 {
			continue
		}

		var steps []*github.TaskStep
		for _, job := range jobs {
			steps = append(steps, job.Steps...)
		}

		stepProfile, err := ghaprofiler.ProfileTaskStep(steps)
		if err != nil {
			log.Fatal(err)
		}
		err = ghaprofiler.SortProfileBy(stepProfile, config.SortBy)
		if err != nil {
			log.Fatal(err)
		}

		// reverse slice
		if config.Reverse {
			reversedStepProfile := make(ghaprofiler.TaskStepProfileResult, len(stepProfile))
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

	var profileFormatterInput ghaprofiler.ProfileInput
	for _, jobName := range formatterInputJobNames {
		result := profileResult[jobName]
		profileFormatterInput = append(profileFormatterInput, &ghaprofiler.ProfileForFormatter{
			Name:    jobName,
			Profile: result,
		})
	}

	ghaprofiler.WriteWithFormat(os.Stdout, profileFormatterInput, config.Format)
}
