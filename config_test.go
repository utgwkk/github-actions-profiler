package ghaprofiler

import "testing"

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
