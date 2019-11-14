package simulator

import (
	"log"
	"time"

	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/network"
)

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

	res := &JobResult{
		JobID:     j.JobID(),
		StartTime: time.Now(),
	}
	if a, err := s.runJob(j); err != nil {
		res.Failed = true
		res.AddError(err)
	} else if a != nil {
		res.AddAction(a)
	}
	res.EndTime = time.Now()

	return res, nil
}

func (s *Simulator) runJob(j Job) (*action.Action, error) {
	params := &JobParams{}
	a, err := j.Execute(params)
	if err != nil {
		return nil, err
	}

	err = s.Network.Push(a)
	if err != nil {
		return nil, err
	}

	totalNodes := len(s.Network.Nodes)
	actions := make(chan *action.Action, totalNodes)
	for i := 0; i < totalNodes; i++ {
		go func(i int) {
			node, err := s.Network.Node(i)
			if err != nil {
				actions <- nil
				return
			}

			for {
				a2, err := node.Action(a.Hash)
				if err == network.ErrActionNotFound {
					continue
				} else if err != nil {
					log.Printf("failed to get action from node #%d: %v", i, err)
				}
				actions <- a2
			}
		}(i)
	}

	total := 0
	for i := 0; i < totalNodes; i++ {
		a := <-actions
		if a == nil {
			log.Printf("got a nil action")
		} else {
			total++
		}
	}

	log.Printf("got %d actions from %d nodes", total, totalNodes)

	return a, nil
}
