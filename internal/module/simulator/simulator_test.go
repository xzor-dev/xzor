package simulator_test

import (
	"errors"
	"testing"
	"time"

	"github.com/xzor-dev/xzor/internal/module/simulator"
)

func TestSimulator(t *testing.T) {
	sim := simulator.New()
	t.Run("Run Empty Jobs", func(t *testing.T) {
		if _, err := sim.Run(nil); err != simulator.ErrNoJobs {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	jobA := &testJob{
		config: &simulator.JobConfig{
			NetworkSize: 2,
		},
		callback: func(p *simulator.JobParams) error {
			return nil
		},
	}
	t.Run("Run Single Job", func(t *testing.T) {
		jobRes := sim.RunJob(jobA)
		if len(jobRes.Errors) > 0 {
			t.Fatalf("job produced %d errors", len(jobRes.Errors))
		}
		if jobRes.TotalExecutions != 2 {
			t.Fatalf("expected 2 executions, got %d", jobRes.TotalExecutions)
		}
	})

	jobB := &testJob{
		config: &simulator.JobConfig{
			NetworkSize: 24,
		},
		callback: func(p *simulator.JobParams) error {
			time.Sleep(time.Millisecond)
			return nil
		},
	}
	t.Run("Two Jobs", func(t *testing.T) {
		runRes, err := sim.Run([]simulator.Job{jobA, jobB})
		if err != nil {
			t.Fatalf("%v", err)
		}
		if runRes.TotalCompleted != 2 {
			t.Fatalf("expected 2 job completions, got %d", runRes.TotalCompleted)
		}
	})

	jobC := &testJob{
		config: &simulator.JobConfig{
			NetworkSize: 128,
		},
		callback: func(p *simulator.JobParams) error {
			sleepTime := time.Millisecond * time.Duration(p.ExecutionIndex)
			time.Sleep(sleepTime)
			return nil
		},
	}
	t.Run("Simulated Lag", func(t *testing.T) {
		jobRes := sim.RunJob(jobC)
		if jobRes.Failed {
			t.Fatalf("job failed with %d errors", len(jobRes.Errors))
		}
	})
}

var _ simulator.Job = &testJob{}

type testJob struct {
	config   *simulator.JobConfig
	callback func(*simulator.JobParams) error
}

func (j *testJob) Config() *simulator.JobConfig {
	return j.config
}

func (j *testJob) Execute(params *simulator.JobParams) error {
	if j.callback == nil {
		return errors.New("no callback function")
	}
	return j.callback(params)
}
