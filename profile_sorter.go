package ghaprofiler

import "sort"

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
