package network

import (
	"bufio"
	"io"
	"net"
	"sync"

	"github.com/xzor-dev/xzor/internal/xzor/action"
)

// Node handles sending messages to and receiving messages from nodes on the network.
type Node struct {
	actionChan  chan *action.Action
	actionMap   map[action.Hash]bool
	connections []net.Conn
	listeners   []net.Listener
	mu          sync.Mutex
}

// NewNode creates a new node instance.
func NewNode() *Node {
	return &Node{
		actionChan: make(chan *action.Action),
		actionMap:  make(map[action.Hash]bool),
	}
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

// Read returns actions as they are recieved by the node.
func (n *Node) Read() (*action.Action, error) {
	return <-n.actionChan, nil
}

// Write sends an action to all registered connections.
func (n *Node) Write(a *action.Action) error {
	if n.actionMap[a.Hash] {
		return action.ErrDuplicateAction
	}
	data, err := a.Encode()
	if err != nil {
		return err
	}

	return n.write(a.Hash, data)
}

func (n *Node) handleIncomingData(data action.EncodedAction) error {
	a, err := data.Decode()
	if err != nil {
		return err
	}
	if n.actionMap[a.Hash] {
		return action.ErrDuplicateAction
	}
	go func() {
		n.actionChan <- a
	}()
	return n.write(a.Hash, data)
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
		encodedAction := action.EncodedAction(data[:len(data)-1])
		n.handleIncomingData(encodedAction)
	}
}

func (n *Node) write(hash action.Hash, encodedAction action.EncodedAction) error {
	n.mu.Lock()
	n.actionMap[hash] = true
	n.mu.Unlock()

	data := append(encodedAction, '\n')
	for _, conn := range n.connections {
		go conn.Write(data)
	}
	return nil
}
