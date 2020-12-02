package ghaprofiler

import (
	"sync"

	"github.com/google/go-github/v32/github"
)

type mapType = map[string][]*github.WorkflowJob

type jobsByJobNameMap struct {
	sync.RWMutex
	m mapType
}

func NewJobsByJobNameMap() *jobsByJobNameMap {
	return &jobsByJobNameMap{
		m: mapType{},
	}
}

func (m *jobsByJobNameMap) Append(jobName string, job *github.WorkflowJob) {
	m.Lock()
	m.m[jobName] = append(m.m[jobName], job)
	m.Unlock()
}

func (m *jobsByJobNameMap) Concat(jobName string, xs []*github.WorkflowJob) {
	m.Lock()
	m.m[jobName] = append(m.m[jobName], xs...)
	m.Unlock()
}

func (m *jobsByJobNameMap) GetJobsByName(jobName string) []*github.WorkflowJob {
	m.RLock()
	defer m.RUnlock()
	return m.m[jobName]
}

func (m *jobsByJobNameMap) Iterate() mapType {
	return m.m
}
