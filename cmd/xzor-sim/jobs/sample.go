package jobs

import (
	"time"

	"github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"
	"github.com/xzor-dev/xzor/internal/xzor/action"
)

func init() {
	Registry.Set(&SampleJobFactory{})
}

const sampleJobID = "sample"

type TestJob struct{}

func (j *TestJob) Config() *simulator.JobConfig {
	return &simulator.JobConfig{
		ExecutionCount: 12,
	}
}

func (j *TestJob) JobID() simulator.JobID {
	return sampleJobID
}

func (j *TestJob) Execute(p *simulator.JobParams) (*action.Action, error) {
	time.Sleep(time.Second * 2)
	return action.New(sampleJobID, "sample-cmd", nil)
}

type SampleJobFactory struct{}

func (j *SampleJobFactory) JobID() simulator.JobID {
	return sampleJobID
}

func (j *SampleJobFactory) NewJob() (simulator.Job, error) {
	return &TestJob{}, nil
}
