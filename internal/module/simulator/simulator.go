package simulator

type Simulator struct{}

func (s *Simulator) Run(jobs []Job) (*JobsResult, error) {
	if len(jobs) == 0 {
		return nil, ErrNoJobs
	}
	res := &JobsResult{}
	jobResChan := make(chan *JobResult, len(jobs))
	for _, j := range jobs {
		go func(j Job) {
			jobRes := s.RunJob(j)
			jobResChan <- jobRes
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
	return res, nil
}

func (s *Simulator) RunJob(j Job) *JobResult {
	c := j.Config()
	res := &JobResult{}
	errChan := make(chan error, c.NetworkSize)
	for i := 0; i < c.NetworkSize; i++ {
		go func(i int) {
			params := &JobParams{
				ExecutionIndex: i,
			}
			err := j.Execute(params)
			errChan <- err
		}(i)
	}
	for i := 0; i < c.NetworkSize; i++ {
		err := <-errChan
		res.TotalExecutions++
		if err != nil {
			res.Failed = true
			res.Errors = append(res.Errors, err)
		}
	}
	return res
}

func New() *Simulator {
	return &Simulator{}
}
