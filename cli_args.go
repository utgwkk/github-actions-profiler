package ghaprofiler

// ProfileConfigCLIArgs is a set of option from command-line arguments
// see DefaultProfileConfig() in config.go for more details
type ProfileConfigCLIArgs struct {
	AccessToken      *string `long:"access_token" description:"Access token for GitHub"`
	Cache            *bool   `long:"cache" description:"Enable disk cache" default-mask:"true"`
	CacheDirectory   *string `long:"cache-dir" description:"Where to store cache data"`
	Concurrency      *int    `long:"concurrency" description:"Concurrency of GitHub API client" default-mask:"2"`
	ConfigPath       *string `long:"config" description:"Path to configuration TOML file"`
	Count            *int    `long:"count" default-mask:"20"`
	Format           *string `long:"format" description:"Output format" default-mask:"table" choice:"table" choice:"json"`
	JobNameRegexp    *string `long:"job-name-regexp" description:"Filter regular expression for a job name"`
	Owner            *string `long:"owner" description:"Repository owner name"`
	Repository       *string `long:"repository" description:"Repository name"`
	Reverse          *bool   `long:"reverse" description:"Reverse the result of sort" default-mask:"false"`
	SortBy           *string `long:"sort" description:"A field name to sort by" default-mask:"number"`
	Verbose          *bool   `long:"verbose" description:"Verbose mode"`
	WorkflowFileName *string `long:"workflow_file" description:"Workflow file name"`
}

func OverrideCLIArgs(tomlConfig *ProfileConfig, cliArgs *ProfileConfigCLIArgs) (newConfig *ProfileConfig) {
	newConfig = &ProfileConfig{}
	if cliArgs.AccessToken != nil {
		newConfig.AccessToken = *cliArgs.AccessToken
	} else {
		newConfig.AccessToken = tomlConfig.AccessToken
	}
	if cliArgs.Cache != nil {
		newConfig.Cache = *cliArgs.Cache
	} else {
		newConfig.Cache = tomlConfig.Cache
	}
	if cliArgs.CacheDirectory != nil {
		newConfig.CacheDirectory = *cliArgs.CacheDirectory
	} else {
		newConfig.CacheDirectory = tomlConfig.CacheDirectory
	}
	if cliArgs.Concurrency != nil {
		newConfig.Concurrency = *cliArgs.Concurrency
	} else {
		newConfig.Concurrency = tomlConfig.Concurrency
	}
	if cliArgs.Count != nil {
		newConfig.Count = *cliArgs.Count
	} else {
		newConfig.Count = tomlConfig.Count
	}
	if cliArgs.Format != nil {
		newConfig.Format = *cliArgs.Format
	} else {
		newConfig.Format = tomlConfig.Format
	}
	if cliArgs.JobNameRegexp != nil {
		newConfig.JobNameRegexp = *cliArgs.JobNameRegexp
	} else {
		newConfig.JobNameRegexp = tomlConfig.JobNameRegexp
	}
	if cliArgs.Owner != nil {
		newConfig.Owner = *cliArgs.Owner
	} else {
		newConfig.Owner = tomlConfig.Owner
	}
	if cliArgs.Repository != nil {
		newConfig.Repository = *cliArgs.Repository
	} else {
		newConfig.Repository = tomlConfig.Repository
	}
	newConfig.Replace = tomlConfig.Replace
	if cliArgs.Reverse != nil {
		newConfig.Reverse = *cliArgs.Reverse
	} else {
		newConfig.Reverse = tomlConfig.Reverse
	}
	if cliArgs.SortBy != nil {
		newConfig.SortBy = *cliArgs.SortBy
	} else {
		newConfig.SortBy = tomlConfig.SortBy
	}
	if cliArgs.Verbose != nil {
		newConfig.Verbose = *cliArgs.Verbose
	} else {
		newConfig.Verbose = tomlConfig.Verbose
	}
	if cliArgs.WorkflowFileName != nil {
		newConfig.WorkflowFileName = *cliArgs.WorkflowFileName
	} else {
		newConfig.WorkflowFileName = tomlConfig.WorkflowFileName
	}
	return
}
