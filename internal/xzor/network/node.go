package network

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"sync"
)

// Node handles sending messages to and receiving messages from nodes on the network.
type Node struct {
	connections []net.Conn
	id          int
	listeners   []net.Listener
	messageChan chan *Message
	messages    map[MessageHash]*Message
	mu          sync.Mutex
}

// AddConnection adds a new remote connection to the node.
func (n *Node) AddConnection(conn net.Conn) {
	n.connections = append(n.connections, conn)
}

// AddListener adds a new local listener.
func (n *Node) AddListener(listener net.Listener) {
	n.listeners = append(n.listeners, listener)
	go n.handleListener(listener)
}

// Connections returns all external connections for the node.
func (n *Node) Connections() []net.Conn {
	return n.connections
}

// Read returns the last message recieved by any of the node's listeners.
func (n *Node) Read() (*Message, error) {
	return <-n.messageChan, nil
}

// Write sends a message to all registered connections.
func (n *Node) Write(msg *Message) error {
	if !n.addMessage(msg) {
		return nil
	}

	return n.write(msg)
}

func (n *Node) addMessage(msg *Message) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.messages == nil {
		n.messages = make(map[MessageHash]*Message)
	}
	if n.messages[msg.Hash] == nil {
		n.messages[msg.Hash] = msg
		return true
	}
	return false
}

func (n *Node) handleIncomingData(data []byte) error {
	msg := &Message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}

	if n.addMessage(msg) {
		go func() {
			n.messageChan <- msg
		}()
		return n.write(msg)
	}

	return nil
}

func (n *Node) handleListener(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			break
		}
		go n.handleListenerConnection(conn)
	}
}

func (n *Node) handleListenerConnection(conn net.Conn) {
	defer conn.Close()
	buffer := bufio.NewReader(conn)
	for {
		data, err := buffer.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}
		n.handleIncomingData(data[:len(data)-1])
	}
}

func (n *Node) write(msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	for _, conn := range n.connections {
		go conn.Write(data)
	}
	return nil
}

var nodeID = 0

// NewNode creates a new node instance.
func NewNode() *Node {
	id := nodeID
	nodeID++
	return &Node{
		id:          id,
		messageChan: make(chan *Message),
		messages:    make(map[MessageHash]*Message),
	}
}
