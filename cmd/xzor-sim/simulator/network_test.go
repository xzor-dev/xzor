package simulator_test

import (
	"fmt"
	"testing"

	"github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"
	"github.com/xzor-dev/xzor/internal/xzor/action"
)

func TestSimulatedNetwork(t *testing.T) {
	networkSize := 32
	connsPerNode := 6

	builder := simulator.NewNetworkLoopBuilder(networkSize, connsPerNode)
	nodes, err := builder.Build()
	if err != nil {
		t.Fatalf("%v", err)
	}

	n, err := simulator.NewNetwork(nodes)
	if err != nil {
		t.Fatalf("%v", err)
	}

	rootNode, err := n.Node(0)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(rootNode.Connections()) != connsPerNode {
		t.Fatalf("expected root node to have %d connections, got %d", connsPerNode, len(rootNode.Connections()))
	}

	actionA, err := action.New("test-mod", "test-cmd", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = rootNode.Write(actionA)
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
			a, err := node.Read()
			if err != nil {
				errors <- err
			} else if a.Hash != actionA.Hash {
				errors <- fmt.Errorf("expected action hash %s from node%d, got %s", actionA.Hash, i, a.Hash)
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
	builder := simulator.NewNetworkWebBuilder(8, 3)
	nodes, err := builder.Build()
	if err != nil {
		t.Fatalf("%v", err)
	}
	n, err := simulator.NewNetwork(nodes)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if node, err := n.Node(0); err != nil {
		t.Fatalf("%v", err)
	} else if len(node.Connections()) != builder.MaxConnections {
		t.Fatalf("expected node #0 to have %d connections, got %d", builder.MaxConnections, len(node.Connections()))
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
		a, err := action.New("test-mod", "test-cmd", nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		err = node5.Write(a)

		node0, err := n.Node(0)
		if err != nil {
			t.Fatalf("%v", err)
		}
		a0, err := node0.Read()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if a0.Hash != a.Hash {
			t.Fatalf("expected action hash %s, got %s", a.Hash, a0.Hash)
		}
	})
}
