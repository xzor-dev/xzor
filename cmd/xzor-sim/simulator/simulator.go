package simulator

import "log"

// Config contains options for the simulator.
type Config struct{}

// Simulator runs single or batches of jobs.
type Simulator struct {
	Config  *Config
	Network *Network
}

// New creates a new simulator instance with the supplied config.
func New(c *Config, n *Network) *Simulator {
	return &Simulator{
		Config:  c,
		Network: n,
	}
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
	log.Printf("executing job %s %d time(s)", j.JobID(), c.ExecutionCount)
	for i := 0; i < c.ExecutionCount; i++ {
		go func(i int) {
			log.Printf("execution #%d of job %s", i+1, j.JobID())
			params := &JobParams{
				ExecutionIndex: i,
			}
			a, err := j.Execute(params)
			if err != nil {
				errChan <- err
				return
			}
			res.AddAction(a)

			log.Printf("pushing action %s to network", a.Hash)
			err = s.Network.Push(a)
			if err != nil {
				log.Printf("failed to push action to network: %v", err)
			}
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
