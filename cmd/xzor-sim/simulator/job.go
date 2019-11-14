package simulator

import (
	"sync"
	"time"

	"github.com/xzor-dev/xzor/internal/xzor/action"
)

// Job is the interface for all simulator jobs.
type Job interface {
	Config() *JobConfig
	Execute(*JobParams) (*action.Action, error)
	JobID() JobID
}

// JobConfig contains configuration values for a job.
type JobConfig struct {
	ExecutionCount    int           // Number of times to execute the job.
	ExecutionInterval time.Duration // Duration between each execution.
}

// JobID is a string used to identify a job.
type JobID string

// JobParams are provided to each job during execution.
type JobParams struct {
	ExecutionIndex int
}

// JobResult contains information for a single job.
type JobResult struct {
	Actions         []*action.Action
	EndTime         time.Time
	Errors          []error
	Failed          bool
	JobID           JobID
	StartTime       time.Time
	TotalExecutions int

	mux sync.Mutex
}

// AddAction adds an action to the job result in a thread-safe way.
func (r *JobResult) AddAction(a *action.Action) {
	r.mux.Lock()
	r.Actions = append(r.Actions, a)
	r.mux.Unlock()
}

// AddError add an error to the job result in a thread-safe way.
func (r *JobResult) AddError(err error) {
	r.mux.Lock()
	r.Errors = append(r.Errors, err)
	r.mux.Unlock()
}

// JobsResult contains information for a group of jobs.
type JobsResult struct {
	Errors         []error
	Jobs           []*JobResult
	TotalCompleted int
	TotalExecuted  int
	TotalFailed    int
}
