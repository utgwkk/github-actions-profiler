package ghaprofiler

import (
	"fmt"
	"sort"
	"strings"
)

var availableSortFields = []string{
	"number",
	"min",
	"max",
	"mean",
	"median",
	"p50",
	"p90",
	"p95",
	"p99",
}

type taskStepProfileSorter struct {
	taskStepProfiles []*TaskStepProfile
	by               func(t1, t2 *TaskStepProfile) bool
}

type taskStepProfileSortBy func(t1, t2 *TaskStepProfile) bool

func (by taskStepProfileSortBy) Sort(taskSteps []*TaskStepProfile) {
	ts := &taskStepProfileSorter{
		taskStepProfiles: taskSteps,
		by:               by,
	}
	sort.Sort(ts)
}

func (ts *taskStepProfileSorter) Len() int {
	return len(ts.taskStepProfiles)
}

func (ts *taskStepProfileSorter) Swap(i, j int) {
	ts.taskStepProfiles[i], ts.taskStepProfiles[j] = ts.taskStepProfiles[j], ts.taskStepProfiles[i]
}

func (ts *taskStepProfileSorter) Less(i, j int) bool {
	return ts.by(ts.taskStepProfiles[i], ts.taskStepProfiles[j])
}

func SortProfileBy(profile TaskStepProfileResult, fieldName string) error {
	var by taskStepProfileSortBy

	switch fieldName {
	case "number":
		by = func(t1, t2 *TaskStepProfile) bool {
			return t1.Number < t2.Number
		}
	case "min":
		by = func(t1, t2 *TaskStepProfile) bool {
			return t1.Min < t2.Min
		}
	case "max":
		by = func(t1, t2 *TaskStepProfile) bool {
			return t1.Max < t2.Max
		}
	case "mean":
		by = func(t1, t2 *TaskStepProfile) bool {
			return t1.Mean < t2.Mean
		}
	case "median":
		by = func(t1, t2 *TaskStepProfile) bool {
			return t1.Median < t2.Median
		}
	case "p50":
		by = func(t1, t2 *TaskStepProfile) bool {
			return t1.Percentile50 < t2.Percentile50
		}
	case "p90":
		by = func(t1, t2 *TaskStepProfile) bool {
			return t1.Percentile90 < t2.Percentile90
		}
	case "p95":
		by = func(t1, t2 *TaskStepProfile) bool {
			return t1.Percentile95 < t2.Percentile95
		}
	case "p99":
		by = func(t1, t2 *TaskStepProfile) bool {
			return t1.Percentile99 < t2.Percentile99
		}
	default:
		return fmt.Errorf("Invalid field: %s", fieldName)
	}
	taskStepProfileSortBy(by).Sort(profile)
	return nil
}

func AvailableSortFieldsForCLI() string {
	return strings.Join(availableSortFields, ", ")
}

func IsValidSortFieldName(fieldName string) bool {
	for _, availableName := range availableSortFields {
		if fieldName == availableName {
			return true
		}
	}
	return false
}
