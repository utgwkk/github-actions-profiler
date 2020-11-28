package ghaprofiler

import "testing"

func Test_OverrideDefaultConfig(t *testing.T) {
	beforeConfig := DefaultProfileConfig()
	cliArgs := &ProfileConfigCLIArgs{
		Verbose: new(bool),
	}
	*cliArgs.Verbose = true

	config := OverrideCLIArgs(beforeConfig, cliArgs)
	if !config.Verbose {
		t.Fatal("Expected --verbose")
	}
}

func Test_OverrideToFalse(t *testing.T) {
	beforeConfig := DefaultProfileConfig()
	beforeConfig.Reverse = true
	beforeConfig.Verbose = true

	cliArgs := &ProfileConfigCLIArgs{
		Reverse: new(bool),
		Verbose: new(bool),
	}
	*cliArgs.Reverse = false
	*cliArgs.Verbose = false

	config := OverrideCLIArgs(beforeConfig, cliArgs)
	if config.Verbose {
		t.Fatal("Unexpected --verbose")
	}
	if config.Reverse {
		t.Fatal("Unexpected --reverse")
	}
}
