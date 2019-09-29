package simulator

type Job interface {
	Config() *JobConfig
	Execute(*JobParams) error
}

type JobConfig struct {
	NetworkSize int
}

type JobParams struct {
	ExecutionIndex int
}

type JobResult struct {
	Failed          bool
	Errors          []error
	TotalExecutions int
}

type JobsResult struct {
	Jobs           []*JobResult
	TotalCompleted int
	TotalExecuted  int
	TotalFailed    int
}
