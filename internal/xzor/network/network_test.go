package network_test

import (
	"net"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/network"
)

func TestBasicNetwork(t *testing.T) {
	nodeA := network.NewNode()
	nodeB := network.NewNode()

	pipeA, pipeB := net.Pipe()
	a, err := action.New("test-mod", "test-cmd", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	nodeA.AddConnection(pipeB)
	nodeB.AddListener(&network.MockListener{
		Connections: []net.Conn{pipeA},
	})

	go func() {
		err := nodeA.Write(a)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()

	actionB, err := nodeB.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if actionB.Hash != a.Hash {
		t.Fatalf("mismatched action hashes: wanted %s, got %s", a.Hash, actionB.Hash)
	}
}

func TestNetworkPropagation(t *testing.T) {
	nodeA := network.NewNode()
	nodeB := network.NewNode()
	nodeC := network.NewNode()

	pipeA1, pipeA2 := net.Pipe()
	pipeB1, pipeB2 := net.Pipe()
	pipeC1, pipeC2 := net.Pipe()

	// create network loop:
	// nodeA -> nodeB
	// nodeB -> nodeC
	// nodeC -> nodeA

	nodeA.AddConnection(pipeB1)
	nodeA.AddListener(&network.MockListener{
		Connections: []net.Conn{pipeA2},
	})

	nodeB.AddConnection(pipeC1)
	nodeB.AddListener(&network.MockListener{
		Connections: []net.Conn{pipeB2},
	})

	nodeC.AddConnection(pipeA1)
	nodeC.AddListener(&network.MockListener{
		Connections: []net.Conn{pipeC2},
	})

	actionA, err := action.New("test-mod", "test-cmd", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	go func() {
		err := nodeA.Write(actionA)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()

	nodeBAction, err := nodeB.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if nodeBAction.Hash != actionA.Hash {
		t.Fatalf("wanted action %s from nodeB, got %s", actionA.Hash, nodeBAction.Hash)
	}

	nodeCAction, err := nodeC.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if nodeCAction.Hash != actionA.Hash {
		t.Fatalf("wanted action %s from nodeC, got %s", actionA.Hash, nodeCAction.Hash)
	}

	actionB, err := action.New("test-mod", "test-cmd2", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	go func() {
		err := nodeC.Write(actionB)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()

	nodeAAction, err := nodeA.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if nodeAAction.Hash != actionB.Hash {
		t.Fatalf("expected action %s from nodeA, got %s", actionB.Hash, nodeAAction.Hash)
	}
}
