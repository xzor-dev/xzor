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
	actionHandler action.Handler
	connections   []net.Conn
	listeners     []net.Listener
	//messageChan   chan EncodedMessage
	//messages      map[MessageHash]bool
	mu sync.Mutex
}

// NewNode creates a new node instance.
func NewNode(actionHandler action.Handler) *Node {
	return &Node{
		actionHandler: actionHandler,
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

// Write sends an action to all registered connections.
func (n *Node) Write(a *action.Action) error {
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

	err = n.actionHandler.HandleAction(a)
	if err != nil {
		return err
	}
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
	data := append(encodedAction, '\n')
	for _, conn := range n.connections {
		go conn.Write(data)
	}
	return nil
}
