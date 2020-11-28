package ghaprofiler

import (
	"context"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

type Client struct {
	githubClient *github.Client
}

type ClientConfig struct {
	AccessToken string
}

func NewClientWithConfig(ctx context.Context, config *ClientConfig) *Client {
	client := &Client{
		githubClient: github.NewClient(nil),
	}

	if config == nil {
		return client
	}

	if config.AccessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.AccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client.githubClient = github.NewClient(tc)
	}

	return client
}

func (c Client) GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*github.WorkflowJob, *github.Response, error) {
	return c.githubClient.Actions.GetWorkflowJobByID(ctx, owner, repo, jobID)
}

func (c Client) ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *github.ListWorkflowJobsOptions) (*github.Jobs, *github.Response, error) {
	return c.githubClient.Actions.ListWorkflowJobs(ctx, owner, repo, runID, opts)
}

func (c Client) ListWorkflowRunsByFileName(ctx context.Context, owner, repo, workflowFileName string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error) {
	return c.githubClient.Actions.ListWorkflowRunsByFileName(ctx, owner, repo, workflowFileName, opts)
}
