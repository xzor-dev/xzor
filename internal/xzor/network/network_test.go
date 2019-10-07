package network_test

import (
	"net"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/network"
)

func TestBasicNetwork(t *testing.T) {
	nodeA := network.NewNode()
	nodeB := network.NewNode()
	pipeA, pipeB := net.Pipe()
	message, err := network.NewMessage("test")
	if err != nil {
		t.Fatalf("%v", err)
	}

	nodeA.AddConnection(pipeB)
	nodeB.AddListener(&network.MockListener{
		Connections: []net.Conn{pipeA},
	})

	go func() {
		err := nodeA.Write(message)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()

	messageB, err := nodeB.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if messageB.Data != message.Data {
		t.Fatalf("expected %s from nodeB, got %s", message.Data, messageB.Data)
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

	nodeBMessage, err := nodeB.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if nodeBMessage.Data != messageA.Data {
		t.Fatalf("wanted %s from nodeB, got %s", messageA.Data, nodeBMessage.Data)
	}

	nodeCMessage, err := nodeC.Read()
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

	nodeAMessage, err := nodeA.Read()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if nodeAMessage.Data != messageB.Data {
		t.Fatalf("expected %s from nodeA, got %s", messageB.Data, nodeAMessage.Data)
	}
}
