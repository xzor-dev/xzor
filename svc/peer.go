package svc

import (
	"bufio"
	"log"
	"net"
)

// Peer represents a connection to another Node.
type Peer struct {
	Address string `json:"address"`
}

// Send sends a message to the Node represented by this Peer.
func (p *Peer) Send(msg []byte) error {
	conn, err := net.Dial("tcp", p.Address)
	if err != nil {
		return err
	}
	// defer conn.Close()

	log.Printf("sending message to peer: %s", msg)
	i, err := bufio.NewWriter(conn).Write(msg)
	if err != nil {
		return err
	}
	log.Printf("sent %d bytes of data", i)
	return nil
}
