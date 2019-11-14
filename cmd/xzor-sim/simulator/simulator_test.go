package simulator_test

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"
	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/common"
)

func TestSimulator(t *testing.T) {
	totalNodes := 12
	connectionsPerNode := 5
	builder := simulator.NewNetworkWebBuilder(totalNodes, connectionsPerNode)
	nodes, err := builder.Build()
	if err != nil {
		t.Fatalf("%v", err)
	}
	nw, err := simulator.NewNetwork(nodes)
	firstNode, err := nw.Node(0)
	if err != nil {
		t.Fatalf("%v", err)
	}
	lastNode, err := nw.Node(totalNodes - 1)
	if err != nil {
		t.Fatalf("%v", err)
	}

	firstNodeActions := make(chan *action.Action, 20)
	lastNodeActions := make(chan *action.Action, 20)

	go func() {
		for {
			a, err := firstNode.Read()
			if err != nil {
				log.Printf("failed to read from first node: %v", err)
				continue
			}
			log.Printf("reading action %s from first node", a.Hash)
			firstNodeActions <- a
		}
	}()
	go func() {
		for {
			a, err := lastNode.Read()
			if err != nil {
				log.Printf("failed to read from last node: %v", err)
				continue
			}
			log.Printf("reading action %s from last node", a.Hash)
			lastNodeActions <- a
		}
	}()

	config := &simulator.Config{}
	sim := simulator.New(config, nw)
	t.Run("Run Empty Jobs", func(t *testing.T) {
		if _, err := sim.Run(nil); err != simulator.ErrNoJobs {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	jobA := &testJob{
		jobID:  "job-a",
		config: &simulator.JobConfig{},
		callback: func(p *simulator.JobParams) (*action.Action, error) {
			hash, err := common.NewRandomHash(10)
			if err != nil {
				return nil, err
			}
			return action.New("job-a", "cmd-a", map[string]interface{}{
				"index": p.ExecutionIndex,
				"hash":  hash,
			})
		},
	}
	t.Run("Run Single Job", func(t *testing.T) {
		jobRes, err := sim.RunJob(jobA)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if len(jobRes.Errors) > 0 {
			log.Printf("job produced %d errors", len(jobRes.Errors))
			for i, err := range jobRes.Errors {
				log.Printf("error #%d: %v", i+1, err)
			}
		}

		a1 := <-lastNodeActions
		found := false
		for _, a2 := range jobRes.Actions {
			if a2.Hash == a1.Hash {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected action %s to be in the job result", a1.Hash)
		}
	})

	jobB := &testJob{
		jobID:  "job-b",
		config: &simulator.JobConfig{},
		callback: func(p *simulator.JobParams) (*action.Action, error) {
			time.Sleep(time.Millisecond)

			hash, err := common.NewRandomHash(10)
			if err != nil {
				return nil, err
			}
			return action.New("job-b", "cmd-b", map[string]interface{}{
				"index": p.ExecutionIndex,
				"hash":  hash,
			})
		},
	}
	t.Run("Two Jobs", func(t *testing.T) {
		runRes, err := sim.Run([]simulator.Job{jobA, jobB})
		if err != nil {
			t.Fatalf("%v", err)
		}
		log.Printf("jobs failed: %d", runRes.TotalFailed)
		log.Printf("jobs completed: %d", runRes.TotalCompleted)
		if len(runRes.Errors) > 0 {
			log.Printf("jobs produced %d errors", len(runRes.Errors))
			for i, err := range runRes.Errors {
				log.Printf("error #%d: %v", i+1, err)
			}
		}
		if runRes.TotalCompleted != 2 {
			t.Fatalf("expected 2 job completions, got %d", runRes.TotalCompleted)
		}

	})
}

var _ simulator.Job = &testJob{}

type testJob struct {
	config   *simulator.JobConfig
	callback func(*simulator.JobParams) (*action.Action, error)
	jobID    simulator.JobID
}

func (j *testJob) Config() *simulator.JobConfig {
	return j.config
}

func (j *testJob) Execute(params *simulator.JobParams) (*action.Action, error) {
	if j.callback == nil {
		return nil, errors.New("no callback function")
	}
	return j.callback(params)
}

func (j *testJob) JobID() simulator.JobID {
	return j.jobID
}
