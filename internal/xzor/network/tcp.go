package network

import (
	"log"
	"net"
)

var _ Connection = &TCPConnection{}

// TCPConnection creates network connections over TCP.
type TCPConnection struct {
	Address string
}

// Connect attempts to establish a connection to a remote TCP server.
func (c *TCPConnection) Connect() (net.Conn, error) {
	log.Printf("attempting to connect to %s", c.Address)

	return net.Dial("tcp", c.Address)
}
