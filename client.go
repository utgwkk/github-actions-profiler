package ghaprofiler

import (
	"context"
	"net/http"

	"github.com/google/go-github/v32/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"golang.org/x/oauth2"
)

var userAgent = "github-actions-profiler (+https://github.com/utgwkk/github-actions-profiler)"

type Client struct {
	githubClient *github.Client
}

type ClientConfig struct {
	AccessToken    string
	Cache          bool
	CacheDirectory string
}

func NewClientWithConfig(ctx context.Context, config *ClientConfig) *Client {
	client := &Client{
		githubClient: github.NewClient(nil),
	}

	if config == nil {
		return client
	}

	var cacheTransport *httpcache.Transport
	var oauth2Client *http.Client
	if config.Cache {
		diskCache := diskcache.New(config.CacheDirectory)
		cacheTransport = httpcache.NewTransport(diskCache)
	}

	if config.AccessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.AccessToken},
		)
		oauth2Client = oauth2.NewClient(ctx, ts)
	}

	if oauth2Client != nil {
		if cacheTransport != nil {
			cacheTransport.Transport = oauth2Client.Transport
			client.githubClient = github.NewClient(cacheTransport.Client())
		} else {
			client.githubClient = github.NewClient(oauth2Client)
		}
	} else {
		if cacheTransport != nil {
			client.githubClient = github.NewClient(cacheTransport.Client())
		} else {
			client.githubClient = github.NewClient(nil)
		}
	}

	client.githubClient.UserAgent = userAgent

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
