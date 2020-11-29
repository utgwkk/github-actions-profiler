package ghaprofiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/pelletier/go-toml"
)

type ProfileConfig struct {
	Owner            string        `toml:"owner"`
	Repository       string        `toml:"repository"`
	WorkflowFileName string        `toml:"workflow-file"`
	Cache            bool          `toml:"cache"`
	CacheDirectory   string        `toml:"cache-directory"`
	Concurrency      int           `toml:"concurrency"`
	Count            int           `toml:"count"`
	AccessToken      string        `toml:"access-token"`
	Format           string        `toml:"format"`
	SortBy           string        `toml:"sort"`
	Reverse          bool          `toml:"reverse"`
	Verbose          bool          `toml:"verbose"`
	JobNameRegexp    string        `toml:"job-name-regexp"`
	Replace          []replaceRule `toml:"replace_rule"`
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

func DefaultProfileConfig() *ProfileConfig {
	return &ProfileConfig{
		Concurrency:    2,
		Count:          20,
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
	if config.Count <= 0 {
		return fmt.Errorf("Count must be a positive integer")
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
