package simulator

import "time"

// Job is the interface for all simulator jobs.
type Job interface {
	Config() *JobConfig
	Execute(*JobParams) error
}

// JobConfig contains configuration values for a job.
type JobConfig struct {
	ExecutionCount    int           // Number of times to execute the job.
	ExecutionInterval time.Duration // Duration between each execution.
}

// JobParams are provided to each job during execution.
type JobParams struct {
	ExecutionIndex int
}

// JobResult contains information for a single job.
type JobResult struct {
	Failed          bool
	Errors          []error
	TotalExecutions int
}

// JobsResult contains information for a group of jobs.
type JobsResult struct {
	Errors         []error
	Jobs           []*JobResult
	TotalCompleted int
	TotalExecuted  int
	TotalFailed    int
}
