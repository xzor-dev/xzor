package simulator

// Config contains options for the simulator.
type Config struct {
	NetworkSize int
}

// Simulator runs single or batches of jobs.
type Simulator struct {
	Config *Config
}

// Run multiple jobs and return a jobs result.
func (s *Simulator) Run(jobs []Job) (*JobsResult, error) {
	if len(jobs) == 0 {
		return nil, ErrNoJobs
	}
	if s.Config == nil {
		return nil, ErrNoConfig
	}

	res := &JobsResult{}
	jobResChan := make(chan *JobResult, len(jobs))
	errChan := make(chan error, len(jobs))
	for _, j := range jobs {
		go func(j Job) {
			jobRes, err := s.RunJob(j)
			jobResChan <- jobRes
			errChan <- err
		}(j)
	}
	for i := 0; i < len(jobs); i++ {
		jobRes := <-jobResChan
		res.Jobs = append(res.Jobs, jobRes)
		if jobRes.Failed {
			res.TotalFailed++
		} else {
			res.TotalCompleted++
		}
	}
	for i := 0; i < len(jobs); i++ {
		err := <-errChan
		if err != nil {
			res.Errors = append(res.Errors, err)
		}
	}
	return res, nil
}

// RunJob runs a single job and returns its result.
func (s *Simulator) RunJob(j Job) (*JobResult, error) {
	if s.Config == nil {
		return nil, ErrNoConfig
	}

	res := &JobResult{}
	c := j.Config()
	if c.ExecutionCount == 0 {
		c.ExecutionCount = 1
	}
	errChan := make(chan error, c.ExecutionCount)
	for i := 0; i < c.ExecutionCount; i++ {
		go func(i int) {
			params := &JobParams{
				ExecutionIndex: i,
			}
			err := j.Execute(params)
			errChan <- err
		}(i)
	}
	for i := 0; i < c.ExecutionCount; i++ {
		err := <-errChan
		res.TotalExecutions++
		if err != nil {
			res.Failed = true
			res.Errors = append(res.Errors, err)
		}
	}
	return res, nil
}

// New creates a new simulator instance with the supplied config.
func New(c *Config) *Simulator {
	return &Simulator{
		Config: c,
	}
}
