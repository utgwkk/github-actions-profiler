package ghaprofiler

import "fmt"

type ProfileConfig struct {
	Owner            string `toml:"owner"`
	Repository       string `toml:"repo"`
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
