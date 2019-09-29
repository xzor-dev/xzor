package job

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/module/simulator"
)

func init() {
	j := &TestJobFactory{}
	Jobs[j.ID()] = j
}

var TestJobID = ID("test-job")

type TestJob struct{}

func (j *TestJob) Config() *simulator.JobConfig {
	return &simulator.JobConfig{
		ExecutionCount: 12,
	}
}

func (j *TestJob) Execute(p *simulator.JobParams) error {
	return errors.New("test job is broken")
}

type TestJobFactory struct{}

func (j *TestJobFactory) ID() ID {
	return TestJobID
}

func (j *TestJobFactory) NewJob() (simulator.Job, error) {
	return &TestJob{}, nil
}
