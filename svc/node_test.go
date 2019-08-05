package svc_test

import (
	"testing"

	"github.com/xzor-dev/xzor/svc"
)

func TestNodeBasics(t *testing.T) {
	n1Addr := ":11111"
	n2Addr := ":11112"

	p1 := &svc.Peer{
		Address: n1Addr,
	}
	p2 := &svc.Peer{
		Address: n2Addr,
	}
	n1 := &svc.Node{
		Address: n1Addr,
		Peers:   []*svc.Peer{p2},
	}
	err := n1.Listen()
	if err != nil {
		t.Fatalf("failed to start node: %v", err)
	}
	n2 := &svc.Node{
		Address: n2Addr,
		Peers:   []*svc.Peer{p1},
	}
	err = n2.Listen()
	if err != nil {
		t.Fatalf("failed to start node: %v", err)
	}

	t.Run("BasicMessage", func(t *testing.T) {
		testMsg := "Hello World"
		err = n1.Write([]byte(testMsg))
		if err != nil {
			t.Fatalf("failed to broadcast message: %v", err)
		}

		go func() {
			msg, err := n2.Read()
			if err != nil {
				t.Fatalf("failed to listen to node: %v", err)
			}
			if string(msg) != testMsg {
				t.Fatalf("unexpected message received: wanted %s, got %s", testMsg, msg)
			}
		}()
	})
}
