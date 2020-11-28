package ghaprofiler

import (
	"regexp"
)

type replaceRule struct {
	reg     *regexp.Regexp
	Regexp  string `toml:"regexp"`
	Replace string `toml:"replace"`
}

func NewReplaceRule(regexpStr, replace string) (*replaceRule, error) {
	regexp, err := regexp.Compile(regexpStr)
	if err != nil {
		return nil, err
	}
	return &replaceRule{
		reg:     regexp,
		Regexp:  regexpStr,
		Replace: replace,
	}, nil
}

func (r *replaceRule) Apply(jobName string) string {
	return r.reg.ReplaceAllLiteralString(jobName, r.Replace)
}
