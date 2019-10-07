package simulator_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/xzor-dev/xzor/internal/module/simulator"
	"github.com/xzor-dev/xzor/internal/xzor/network"
)

func TestSimulator(t *testing.T) {
	config := &simulator.Config{
		TotalNodes: 2,
	}
	sim := simulator.New(config, nil)
	t.Run("Run Empty Jobs", func(t *testing.T) {
		if _, err := sim.Run(nil); err != simulator.ErrNoJobs {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	jobA := &testJob{
		config: &simulator.JobConfig{
			ExecutionCount: 2,
		},
		callback: func(p *simulator.JobParams) error {
			return nil
		},
	}
	t.Run("Run Single Job", func(t *testing.T) {
		jobRes, err := sim.RunJob(jobA)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if len(jobRes.Errors) > 0 {
			t.Fatalf("job produced %d errors", len(jobRes.Errors))
		}
		if jobRes.TotalExecutions != 2 {
			t.Fatalf("expected 2 executions, got %d", jobRes.TotalExecutions)
		}
	})

	jobB := &testJob{
		config: &simulator.JobConfig{},
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
}

func TestSimulatedNetwork(t *testing.T) {
	networkSize := 32
	connsPerNode := 6
	c := &simulator.Config{
		ConnectionsPerNode: connsPerNode,
		TotalNodes:         networkSize,
	}
	n, err := simulator.NewNetwork(c, &simulator.NetworkLoopBuilder{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(n.Nodes) != networkSize {
		t.Fatalf("expected network to have %d nodes, got %d", networkSize, len(n.Nodes))
	}
	rootNode, err := n.Node(0)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(rootNode.Connections()) != connsPerNode {
		t.Fatalf("expected root node to have %d connections, got %d", connsPerNode, len(rootNode.Connections()))
	}

	messageA, err := network.NewMessage("hello")
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = rootNode.Write(messageA)
	if err != nil {
		t.Fatalf("%v", err)
	}

	errors := make(chan error, networkSize-1)
	for i := 1; i < networkSize; i++ {
		go func(i int) {
			node, err := n.Node(i)
			if err != nil {
				t.Fatalf("%v", err)
			}
			msg, err := node.Read()
			if err != nil {
				errors <- err
			} else if msg.Data != messageA.Data {
				errors <- fmt.Errorf("expected %s from node%d, got %s", messageA.Data, i, msg.Data)
			} else {
				errors <- nil
			}
		}(i)
	}
	for i := 1; i < networkSize; i++ {
		err := <-errors
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
}

func TestNetworkWeb(t *testing.T) {
	c := &simulator.Config{
		ConnectionsPerNode: 3,
		TotalNodes:         8,
	}
	n, err := simulator.NewNetwork(c, &simulator.NetworkWebBuilder{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	if node, err := n.Node(0); err != nil {
		t.Fatalf("%v", err)
	} else if len(node.Connections()) != c.ConnectionsPerNode {
		t.Fatalf("expected node #0 to have %d connections, got %d", c.ConnectionsPerNode, len(node.Connections()))
	}

	if node, err := n.Node(1); err != nil {
		t.Fatalf("%v", err)
	} else if len(node.Connections()) != 4 {
		t.Fatalf("expected node #1 to have %d connections, got %d", 4, len(node.Connections()))
	}

	if node, err := n.Node(2); err != nil {
		t.Fatalf("%v", err)
	} else if len(node.Connections()) != 2 {
		t.Fatalf("expected node #2 to have %d connections, got %d", 2, len(node.Connections()))
	}

	t.Run("Propagate Backwards", func(t *testing.T) {
		node5, err := n.Node(5)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if len(node5.Connections()) != 1 {
			t.Fatalf("expected node #5 to have %d connections, got %d", 1, len(node5.Connections()))
		}
		message, err := network.NewMessage("hello")
		if err != nil {
			t.Fatalf("%v", err)
		}
		err = node5.Write(message)

		node0, err := n.Node(0)
		if err != nil {
			t.Fatalf("%v", err)
		}
		m0, err := node0.Read()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if m0.Data != message.Data {
			t.Fatalf("expected %s from message data, got %s", message.Data, m0.Data)
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
