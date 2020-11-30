package ghaprofiler

import (
	"reflect"
	"testing"
)

type ReplaceRuleWithoutRegex struct {
	Regexp  string
	Replace string
}

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
		NumberOfJob:      100,
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

func TestParseCLIArgs(t *testing.T) {
	config := DefaultProfileConfig()
	config.ConfigPath = func(path string) {
		var err error
		configFromTOML, err := LoadConfigFromTOML(path)
		if err != nil {
			t.Fatal(err)
		}
		config = configFromTOML
	}
	parser := NewCLIParser(config)
	args := []string{
		"--config",
		"fixtures/valid-config.toml",
	}
	args, err := parser.ParseArgs(args)
	if err != nil {
		t.Fatal(err)
	}

	expectedConfig := &ProfileConfig{
		AccessToken:      "YOUR_ACCESS_TOKEN",
		Concurrency:      2,
		Cache:            true,
		CacheDirectory:   "/tmp/cache",
		NumberOfJob:      100,
		Format:           "table",
		Owner:            "utgwkk",
		Repository:       "Twitter-Text",
		SortBy:           "number",
		Verbose:          false,
		Reverse:          false,
		WorkflowFileName: "ci.yml",
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Fatalf("Loaded config is not correct\ngot: %#v\nwant: %#v", config, expectedConfig)
	}
}

func TestParseCLIArgsLoadReplaceRule(t *testing.T) {
	config := DefaultProfileConfig()
	config.ConfigPath = func(path string) {
		var err error
		configFromTOML, err := LoadConfigFromTOML(path)
		if err != nil {
			t.Fatal(err)
		}
		config = configFromTOML
	}
	parser := NewCLIParser(config)
	args := []string{
		"--config",
		"fixtures/valid-config-with-replace-rule.toml",
	}
	args, err := parser.ParseArgs(args)
	if err != nil {
		t.Fatal(err)
	}

	var replaceRules []*ReplaceRuleWithoutRegex
	for _, rule := range config.Replace {
		replaceRules = append(replaceRules, &ReplaceRuleWithoutRegex{
			Regexp:  rule.Regexp,
			Replace: rule.Replace,
		})
	}
	expectedReplaceRules := []*ReplaceRuleWithoutRegex{
		{
			Regexp:  `Perl 5\.[0-9]+`,
			Replace: "Perl",
		},
	}

	if !reflect.DeepEqual(replaceRules, expectedReplaceRules) {
		t.Fatalf("Loaded replace rules mismatched\ngot: %#v\nwant: %#v", replaceRules, expectedReplaceRules)
	}
}
