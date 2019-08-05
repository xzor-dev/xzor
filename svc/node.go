package svc

import (
	"bufio"
	"errors"
	"log"
	"net"
)

const (
	// MessageDelimiter indicates the end of a message.
	MessageDelimiter = '\n'
)

// Node holds data on a local connection.
type Node struct {
	Address  string       `json:"address"`
	Listener net.Listener `json:"-"`
	Peers    []*Peer      `json:"peers"`

	lastMessage chan []byte
	messages    [][]byte
}

// Close closes the TCP listen server.
func (n *Node) Close() error {
	if n.Listener == nil {
		return errors.New("node has not started listening")
	}
	return n.Listener.Close()
}

// Listen initializes a TCP listener using the node's address.
func (n *Node) Listen() error {
	if n.Listener != nil {
		return errors.New("node is already listening")
	}

	l, err := net.Listen("tcp", n.Address)
	if err != nil {
		return err
	}
	n.Listener = l
	n.lastMessage = make(chan []byte, 1)
	n.messages = make([][]byte, 0)
	go n.listen()
	return nil
}

// Read incoming messages.
// Operation is blocked until a message is recieved.
func (n *Node) Read() ([]byte, error) {
	if n.Listener == nil {
		return nil, errors.New("node has not started listening")
	}
	return <-n.lastMessage, nil
}

// Write arbitrary data to the node and broadcast to all connected peers.
func (n *Node) Write(msg []byte) error {
	if n.Peers == nil {
		return errors.New("no peers have been assigned to this node")
	}

	msgWithDelim := append(msg, MessageDelimiter)
	for _, p := range n.Peers {
		err := p.Send(msgWithDelim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) handleConnection(c net.Conn) {
	log.Println("handling new connection")

	for {
		log.Printf("waiting for message")
		msg, err := bufio.NewReader(c).ReadBytes(MessageDelimiter)
		if err != nil {
			log.Printf("got error while reading connection: %v", err)
			n.handleError(err)
			return
		}
		n.handleMessage(msg)
	}
}

func (n *Node) handleError(err error) {
	msg := []byte(err.Error())
	n.handleMessage(msg)
}

func (n *Node) handleMessage(msg []byte) {
	log.Printf("handling new message: %s", msg)
	n.messages = append(n.messages, msg)
	n.lastMessage <- msg
}

func (n *Node) listen() {
	for {
		c, err := n.Listener.Accept()
		if err != nil {
			log.Printf("got error from incoming connection: %v", err)
			n.handleError(err)
			continue
		}
		go n.handleConnection(c)
	}
}
