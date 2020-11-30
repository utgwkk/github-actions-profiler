package ghaprofiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/jessevdk/go-flags"
	"github.com/pelletier/go-toml"
)

type ProfileConfig struct {
	ConfigPath func(string) `long:"config" description:"Path to configuration TOML file"`

	AccessToken      string `toml:"access-token" long:"access-token" description:"Access token for GitHub" env:"GITHUB_ACTIONS_PROFILER_TOKEN"`
	Cache            bool   `toml:"cache" long:"cache" description:"Enable disk cache" default-mask:"true"`
	CacheDirectory   string `toml:"cache-directory" long:"cache-dir" description:"Where to store cache data"`
	Concurrency      int    `toml:"concurrency" long:"concurrency" short:"j" description:"Concurrency of GitHub API client" default-mask:"2"`
	Format           string `toml:"format" long:"number-of-job" short:"n" description:"The number of job to analyze" default-mask:"20"`
	JobNameRegexp    string `toml:"job-name-regexp" long:"format" short:"f" description:"Output format" default-mask:"table" choice:"table" choice:"json" choice:"tsv" choice:"markdown"`
	NumberOfJob      int    `toml:"number-of-job" long:"job-name-regexp" description:"Filter regular expression for a job name"`
	Owner            string `toml:"owner" long:"owner" description:"Repository owner name"`
	Repository       string `toml:"repository" long:"repository" description:"Repository name"`
	Reverse          bool   `toml:"reverse" long:"reverse" short:"r" description:"Reverse the result of sort" default-mask:"false"`
	SortBy           string `toml:"sort" long:"sort" short:"s" description:"A field name to sort by" default-mask:"number"`
	Verbose          bool   `toml:"verbose" long:"verbose" description:"Verbose mode"`
	WorkflowFileName string `toml:"workflow-file" long:"workflow-file" description:"Workflow file name"`

	// Not in CLI arguments
	Replace []replaceRule `toml:"replace_rule"`
}

var defaultCacheDirectoryName = "github-actions-profiler-httpcache"

func defaultCacheDirectoryPath() string {
	userCacheDir, err := os.UserCacheDir()
	if err == nil {
		return path.Join(userCacheDir, defaultCacheDirectoryName)
	}

	// fallback to temporary directory
	return path.Join(os.TempDir(), defaultCacheDirectoryName)
}

func NewCLIParser(dst *ProfileConfig) *flags.Parser {
	return flags.NewParser(dst, flags.Default)
}

func DefaultProfileConfig() *ProfileConfig {
	return &ProfileConfig{
		Concurrency:    2,
		NumberOfJob:    20,
		Cache:          true,
		CacheDirectory: defaultCacheDirectoryPath(),
		Format:         "table",
		SortBy:         "number",
	}
}

func (config ProfileConfig) Validate() error {
	if config.Owner == "" {
		return fmt.Errorf("Repository owner name required")
	}
	if config.Repository == "" {
		return fmt.Errorf("Repository name required")
	}
	if config.WorkflowFileName == "" {
		return fmt.Errorf("Workflow file name required")
	}
	if config.Concurrency <= 0 {
		return fmt.Errorf("Concurrency must be a positive integer")
	}
	if config.NumberOfJob <= 0 {
		return fmt.Errorf("NumberOfJob must be a positive integer")
	}
	if !IsValidFormatName(config.Format) {
		return fmt.Errorf("Invalid format: %s", config.Format)
	}
	if !IsValidSortFieldName(config.SortBy) {
		return fmt.Errorf("Invalid sort field name: %s", config.SortBy)
	}
	if _, err := regexp.Compile(config.JobNameRegexp); err != nil {
		return fmt.Errorf("Invalid regular expression: %v", err)
	}
	if config.Cache && config.CacheDirectory == "" {
		return fmt.Errorf("Cache enabled but no cache directory passed")
	}

	return nil
}

func LoadConfigFromTOML(filename string) (*ProfileConfig, error) {
	config := DefaultProfileConfig()

	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(p, config)
	if err != nil {
		return nil, err
	}

	for i, rule := range config.Replace {
		newRule, err := NewReplaceRule(rule.Regexp, rule.Replace)
		if err != nil {
			return nil, err
		}
		config.Replace[i] = *newRule
	}

	return config, nil
}

func (c ProfileConfig) Dump() string {
	var dump string
	dump += fmt.Sprintf("concurrency=%v\n", c.Concurrency)
	dump += fmt.Sprintf("number-of-job=%v\n", c.NumberOfJob)
	dump += fmt.Sprintf("format=%v\n", c.Format)
	dump += fmt.Sprintf("job-name-regexp=%v\n", c.JobNameRegexp)
	dump += fmt.Sprintf("owner=%v\n", c.Owner)
	dump += fmt.Sprintf("repo=%v\n", c.Repository)
	dump += fmt.Sprintf("reverse=%v\n", c.Reverse)
	dump += fmt.Sprintf("sort=%v\n", c.SortBy)
	// We don't write out token
	if c.AccessToken == "" {
		dump += "access token not set\n"
	} else {
		dump += "access token set\n"
	}
	dump += fmt.Sprintf("workflow-file=%v\n", c.WorkflowFileName)
	dump += fmt.Sprintf("replace=%#v\n", c.Replace)
	dump += fmt.Sprintf("cache=%v\n", c.Cache)
	dump += fmt.Sprintf("cache-directory=%v\n", c.CacheDirectory)
	return dump
}
