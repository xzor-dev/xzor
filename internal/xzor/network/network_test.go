package network_test

import (
	"log"
	"net"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/network"
)

func TestNetwork(t *testing.T) {
	dataChan := make(chan []byte)
	dataHandler := &network.MockDataHandler{
		Handler: func(data []byte) error {
			log.Println("handling data")
			dataChan <- data
			return nil
		},
	}
	nodeA := &network.Node{
		DataHandler: dataHandler,
	}

	connA1, connA2 := net.Pipe()
	listenerA := &network.MockListener{
		Conn: connA1,
	}
	nodeA.AddListener(listenerA)

	connB1, connB2 := net.Pipe()
	nodeA.AddConnection(&network.MockConnection{
		Conn: connB1,
	})

	err := nodeA.Start()
	if err != nil {
		t.Fatalf("%v", err)
	}

	go func() {
		err := <-nodeA.Errors
		t.Fatalf("%v", err)
	}()

	msg := "hello"
	msgBytes := append([]byte(msg), '\n')
	_, err = connA2.Write(msgBytes)
	if err != nil {
		t.Fatalf("%v", err)
	}

	lastMsg := <-dataChan
	if string(lastMsg) != msg {
		t.Fatalf("unexpected message received: wanted %s, got %s", msg, lastMsg)
	}

	remoteMsg := make([]byte, len(msg))
	_, err = connB2.Read(remoteMsg)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if string(remoteMsg) != msg {
		t.Fatalf("unexpected message received: wanted %s, got %s", msg, remoteMsg)
	}
}
