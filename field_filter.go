package ghaprofiler

import "strings"

type filterFields []string

type filterFieldFunc func(fieldName string) bool

func only(onlyFields filterFields) filterFieldFunc {
	return func(fieldName string) bool {
		for _, onlyField := range onlyFields {
			if fieldName == onlyField {
				return true
			}
		}
		return false
	}
}

func exclude(excludedFields filterFields) filterFieldFunc {
	return func(fieldName string) bool {
		for _, excludedField := range excludedFields {
			if fieldName == excludedField {
				return false
			}
		}
		return true
	}
}

func excludePercentile() filterFieldFunc {
	return func(fieldName string) bool {
		return !strings.HasPrefix(fieldName, "p")
	}
}
