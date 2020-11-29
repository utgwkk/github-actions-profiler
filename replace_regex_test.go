package ghaprofiler

import (
	"testing"
)

func Test_Apply(t *testing.T) {
	rule, err := NewReplaceRule(`Perl 5\.[0-9]{1,2}`, "Perl")
	if err != nil {
		t.Fatal(err)
	}
	jobName := "Perl 5.32"
	normalizedJobName := rule.Apply(jobName)
	if normalizedJobName != "Perl" {
		t.Fatalf("normalizedJobName does not match\nexpected: %#v\ngot%#v", "Perl", normalizedJobName)
	}
}
