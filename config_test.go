package ghaprofiler

import (
	"reflect"
	"testing"
)

func Test_Validate(t *testing.T) {
	config, err := LoadConfigFromTOML("fixtures/valid-config.toml")
	if err != nil {
		t.Fatal(err)
	}

	err = config.Validate()
	if err != nil {
		t.Fatalf("Validation failed for ghaprofiler.DefaultProfileConfig(): %v", err)
	}
}

func Test_LoadFromTOML(t *testing.T) {
	expectedConfig := &ProfileConfig{
		AccessToken:      "YOUR_ACCESS_TOKEN",
		Concurrency:      2,
		Cache:            true,
		CacheDirectory:   "/tmp/cache",
		Count:            100,
		Format:           "table",
		Owner:            "utgwkk",
		Repository:       "Twitter-Text",
		SortBy:           "number",
		WorkflowFileName: "ci.yml",
	}

	config, err := LoadConfigFromTOML("fixtures/valid-config.toml")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Fatalf("Loaded config is not correct\ngot: %#v\nwant: %#v", config, expectedConfig)
	}
}
