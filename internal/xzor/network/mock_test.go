package network_test

import (
	"net"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/network"
)

func TestMockListener(t *testing.T) {
	connCount := 2
	conns := make([]net.Conn, connCount)
	errors := make(chan error, connCount)
	message := "hello"
	ln := &network.MockListener{}

	for i := 0; i < connCount; i++ {
		cA, cB := net.Pipe()
		ln.AddConnection(cA)
		conns[i] = cB
		go func() {
			conn, err := ln.Accept()
			if err != nil {
				errors <- err
				return
			}
			msg := make([]byte, len(message))
			_, err = conn.Read(msg)
			if err != nil {
				errors <- err
				return
			}
			errors <- nil
		}()
		_, err := cB.Write([]byte(message))
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
	for i := 0; i < connCount; i++ {
		err := <-errors
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
}
