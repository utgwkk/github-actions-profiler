package ghaprofiler

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
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

func (dst *ProfileConfig) OverrideConfig(src *ProfileConfig) {
	dst.AccessToken = src.AccessToken
	dst.Count = src.Count
	dst.Format = src.Format
	dst.Owner = src.Owner
	dst.Repository = src.Repository
	dst.Reverse = src.Reverse
	dst.SortBy = src.SortBy
	dst.Verbose = src.Verbose
	dst.WorkflowFileName = src.WorkflowFileName
}
