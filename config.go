package ghaprofiler

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/pelletier/go-toml"
)

type ProfileConfig struct {
	Owner            string `toml:"owner"`
	Repository       string `toml:"repository"`
	WorkflowFileName string `toml:"workflow_file"`
	Count            int    `toml:"count"`
	AccessToken      string `toml:"access_token"`
	Format           string `toml:"format"`
	SortBy           string `toml:"sort"`
	Reverse          bool   `toml:"reverse"`
	Verbose          bool   `toml:"verbose"`
	JobNameRegexp    string `toml:"job_name_regexp"`
}

func DefaultProfileConfig() *ProfileConfig {
	return &ProfileConfig{
		Count:  20,
		Format: "table",
		SortBy: "number",
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

	return nil
}

func LoadConfigFromTOML(filename string) (*ProfileConfig, error) {
	config := &ProfileConfig{}

	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(p, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
