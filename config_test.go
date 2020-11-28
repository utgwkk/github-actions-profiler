package ghaprofiler

import "testing"

func Test_DefaultProfileConfigPassesValidation(t *testing.T) {
	config := DefaultProfileConfig()
	err := config.Validate()
	if err != nil {
		t.Fatalf("Validation failed for ghaprofiler.DefaultProfileConfig(): %v", err)
	}
}
