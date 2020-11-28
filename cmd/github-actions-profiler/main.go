package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v32/github"
	ghaprofiler "github.com/utgwkk/github-actions-profiler"
)

var owner string
var repo string
var workflowFileName string
var count int
var accessToken string
var format string
var sortBy string
var reverse bool
var verbose bool

const accessTokenEnvVariableName = "GITHUB_ACTIONS_PROFILER_TOKEN"

func init() {
	flag.StringVar(&owner, "owner", "", "Repository owner name")
	flag.StringVar(&repo, "repo", "", "Repository name")
	flag.StringVar(&workflowFileName, "workflow_file", "", "Workflow file name")
	flag.StringVar(&accessToken, "token", "", "Access token. You can pass it with "+accessTokenEnvVariableName+" environment variable")
	flag.IntVar(&count, "count", 20, "Count")
	flag.StringVar(&format, "format", "table", "Output format. Supported formats are: "+ghaprofiler.AvailableFormats())
	flag.StringVar(&sortBy, "sort", "number", "A filed name to sort by. Supported values are"+ghaprofiler.AvailableSortFieldsForCLI())
	flag.BoolVar(&reverse, "reverse", false, "Reverse the result of sort")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode")
}

func validateFlags() error {
	if owner == "" {
		return fmt.Errorf("Repository owner name required")
	}
	if repo == "" {
		return fmt.Errorf("Repository name required")
	}
	if workflowFileName == "" {
		return fmt.Errorf("Workflow file name required")
	}
	if count <= 0 {
		return fmt.Errorf("Count must be a positive integer")
	}
	if !ghaprofiler.IsValidFormatName(format) {
		return fmt.Errorf("Invalid format: %s", format)
	}
	if !ghaprofiler.IsValidSortFieldName(sortBy) {
		return fmt.Errorf("Invalid sort field name: %s", sortBy)
	}

	return nil
}

func main() {
	ctx := context.Background()
	flag.Parse()

	if verbose {
		log.Printf("count=%v\n", count)
		log.Printf("format=%v\n", format)
		log.Printf("owner=%v\n", owner)
		log.Printf("repo=%v\n", repo)
		log.Printf("reverse=%v\n", reverse)
		log.Printf("sort=%v\n", sortBy)
		// We don't write out token
		log.Printf("workflow_file=%v\n", workflowFileName)
	}

	if err := validateFlags(); err != nil {
		log.Fatal(err)
	}

	if accessToken == "" {
		accessTokenFromEnv := os.Getenv(accessTokenEnvVariableName)
		if accessTokenFromEnv != "" {
			accessToken = accessTokenFromEnv
		}
	}

	client := ghaprofiler.NewClientWithConfig(ctx, &ghaprofiler.ClientConfig{
		AccessToken: accessToken,
	})

	listWorkflowRunsOpts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			PerPage: count,
		},
	}

	if verbose {
		log.Println("ListWorkflowRunsByFileName start")
	}
	workflowRuns, _, err := client.ListWorkflowRunsByFileName(ctx, owner, repo, workflowFileName, listWorkflowRunsOpts)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Println("ListWorkflowRunsByFileName finish")
	}

	jobsByJobName := make(map[string][]*github.WorkflowJob)

	for _, run := range workflowRuns.WorkflowRuns {
		if verbose {
			log.Printf("ListWorkflowJobs start: run_id=%d", *run.ID)
		}
		jobs, _, err := client.ListWorkflowJobs(ctx, owner, repo, *run.ID, nil)
		if err != nil {
			log.Fatal(err)
		}
		if verbose {
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
		err = ghaprofiler.SortProfileBy(stepProfile, sortBy)
		if err != nil {
			log.Fatal(err)
		}

		// reverse slice
		if reverse {
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

	ghaprofiler.WriteWithFormat(os.Stdout, profileFormatterInput, format)
}
