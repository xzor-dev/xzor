package network_test

import (
	"net"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/network"
)

func TestBasicNetwork(t *testing.T) {
	handlerA := newTestActionHandler()
	nodeA := network.NewNode(handlerA)

	handlerB := newTestActionHandler()
	nodeB := network.NewNode(handlerB)

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

	actionB := <-handlerB.actions
	if actionB.Hash != a.Hash {
		t.Fatalf("mismatched action hashes: wanted %s, got %s", a.Hash, actionB.Hash)
	}
}

/*
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

	messageA, err := network.NewMessage("test")
	if err != nil {
		t.Fatalf("%v", err)
	}
	go func() {
		err := nodeA.Write(messageA)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()

	nodeBData, err := nodeB.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	nodeBMessage, err := nodeBData.Decode(nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if nodeBMessage.Data != messageA.Data {
		t.Fatalf("wanted %s from nodeB, got %s", messageA.Data, nodeBMessage.Data)
	}

	nodeCData, err := nodeC.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	nodeCMessage, err := nodeCData.Decode(nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if nodeCMessage.Data != messageA.Data {
		t.Fatalf("wanted %s from nodeC, got %s", messageA.Data, nodeCMessage.Data)
	}

	messageB, err := network.NewMessage("test2")
	if err != nil {
		t.Fatalf("%v", err)
	}
	go func() {
		err := nodeC.Write(messageB)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()

	nodeAData, err := nodeA.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	nodeAMessage, err := nodeAData.Decode(nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if nodeAMessage.Data != messageB.Data {
		t.Fatalf("expected %s from nodeA, got %s", messageB.Data, nodeAMessage.Data)
	}
}
*/

var _ action.Handler = &testActionHandler{}

type testActionHandler struct {
	actions chan *action.Action
}

func newTestActionHandler() *testActionHandler {
	return &testActionHandler{
		actions: make(chan *action.Action),
	}
}

func (h *testActionHandler) HandleAction(a *action.Action) error {
	h.actions <- a
	return nil
}
