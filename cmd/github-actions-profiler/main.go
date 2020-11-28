package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/go-github/v32/github"
	ghaprofiler "github.com/utgwkk/github-actions-profiler"
)

var config *ghaprofiler.ProfileConfig = ghaprofiler.DefaultProfileConfig()
var configFromArgs *ghaprofiler.ProfileConfig = ghaprofiler.DefaultProfileConfig()
var configTomlPath string

const accessTokenEnvVariableName = "GITHUB_ACTIONS_PROFILER_TOKEN"

func init() {
	flag.StringVar(&configFromArgs.Owner, "owner", "", "Repository owner name")
	flag.StringVar(&configFromArgs.Repository, "repository", "", "Repository name")
	flag.StringVar(&configFromArgs.WorkflowFileName, "workflow_file", "", "Workflow file name")
	flag.StringVar(&configFromArgs.AccessToken, "access_token", "", "Access token. You can pass it with "+accessTokenEnvVariableName+" environment variable")
	flag.IntVar(&configFromArgs.Count, "count", 20, "Count")
	flag.StringVar(&configFromArgs.Format, "format", "table", "Output format. Supported formats are: "+ghaprofiler.AvailableFormatsForCLI())
	flag.StringVar(&configFromArgs.SortBy, "sort", "number", "A filed name to sort by. Supported values are"+ghaprofiler.AvailableSortFieldsForCLI())
	flag.BoolVar(&configFromArgs.Reverse, "reverse", false, "Reverse the result of sort")
	flag.BoolVar(&configFromArgs.Verbose, "verbose", false, "Verbose mode")
	flag.StringVar(&configTomlPath, "config", "", "Path to configuration TOML file. Note that settings in command-line arguments are overwritten with TOML")
}

func main() {
	ctx := context.Background()
	flag.Parse()

	if configTomlPath != "" {
		configFromTOML, err := ghaprofiler.LoadConfigFromTOML(configTomlPath)
		if err != nil {
			log.Fatalf("Failed to load %s: %v", configTomlPath, err)
		}
		// TODO: Override config with CLI arguments when they are given
		config = configFromTOML
	} else {
		config = configFromArgs
	}

	if configFromArgs.Verbose {
		config.Verbose = configFromArgs.Verbose
	}

	if config.Verbose {
		log.Printf("config=%v\n", configTomlPath)
		log.Printf("count=%v\n", config.Count)
		log.Printf("format=%v\n", config.Format)
		log.Printf("owner=%v\n", config.Owner)
		log.Printf("repo=%v\n", config.Repository)
		log.Printf("reverse=%v\n", config.Reverse)
		log.Printf("sort=%v\n", config.SortBy)
		// We don't write out token
		log.Printf("workflow_file=%v\n", config.WorkflowFileName)
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

	client := ghaprofiler.NewClientWithConfig(ctx, &ghaprofiler.ClientConfig{
		AccessToken: config.AccessToken,
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

	jobsByJobName := make(map[string][]*github.WorkflowJob)

	for _, run := range workflowRuns.WorkflowRuns {
		if config.Verbose {
			log.Printf("ListWorkflowJobs start: run_id=%d", *run.ID)
		}
		jobs, _, err := client.ListWorkflowJobs(ctx, config.Owner, config.Repository, *run.ID, nil)
		if err != nil {
			log.Fatal(err)
		}
		if config.Verbose {
			log.Printf("ListWorkflowJobs finish: run_id=%d", *run.ID)
		}

		for _, job := range jobs.Jobs {
			jobsByJobName[*job.Name] = append(jobsByJobName[*job.Name], job)
		}
	}

	profileResult := make(map[string][]*ghaprofiler.TaskStepProfile)

	for jobName, jobs := range jobsByJobName {
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

	var profileFormatterInput ghaprofiler.ProfileInput
	for jobName, result := range profileResult {
		profileFormatterInput = append(profileFormatterInput, &ghaprofiler.ProfileForFormatter{
			Name:    jobName,
			Profile: result,
		})
	}

	ghaprofiler.WriteWithFormat(os.Stdout, profileFormatterInput, config.Format)
}
