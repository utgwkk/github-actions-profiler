package ghaprofiler

import "strings"

type filterFields []string

type filterFieldFunc func(fieldName string) bool

func Only(onlyFields filterFields) filterFieldFunc {
	return func(fieldName string) bool {
		for _, onlyField := range onlyFields {
			if fieldName == onlyField {
				return true
			}
		}
		return false
	}
}

func Exclude(excludedFields filterFields) filterFieldFunc {
	return func(fieldName string) bool {
		for _, excludedField := range excludedFields {
			if fieldName == excludedField {
				return false
			}
		}
		return true
	}
}

func ExcludePercentile() filterFieldFunc {
	return func(fieldName string) bool {
		return !strings.HasPrefix(fieldName, "p")
	}
}

func ShowAll() filterFieldFunc {
	return func(_ string) bool {
		return true
	}
}
