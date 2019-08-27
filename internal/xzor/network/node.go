package network

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

// Node controls all local operations.
type Node struct {
	DataHandler DataHandler
	Errors      chan error

	connections         []Connection
	inboundConnections  []net.Conn
	listeners           []Listener
	outboundConnections []net.Conn
	quit                chan bool
}

// AddConnection adds a new remote connection to the node.
func (n *Node) AddConnection(conn Connection) {
	if n.connections == nil {
		n.connections = make([]Connection, 0)
	}
	n.connections = append(n.connections, conn)
}

// AddListener adds a new local listener.
func (n *Node) AddListener(listener Listener) {
	if n.listeners == nil {
		n.listeners = make([]Listener, 0)
	}
	n.listeners = append(n.listeners, listener)
}

// Start all components within the node.
func (n *Node) Start() error {
	log.Println("starting node")

	if n.DataHandler == nil {
		return errors.New("no DataHandler provided to the node")
	}

	n.Errors = make(chan error)
	n.initConnections()
	n.initListeners()

	return nil
}

func (n *Node) handleData(data []byte) {
	if string(data[0:4]) == "quit" {
		log.Println("quitting")
		n.quit <- true
		return
	}

	err := n.DataHandler.HandleData(data)
	if err != nil {
		log.Printf("failed to handle data: %v", err)
		n.Errors <- err
		return
	}

	for _, conn := range n.outboundConnections {
		go conn.Write(data)
	}
}

func (n *Node) handleInboundConnection(conn net.Conn) {
	n.inboundConnections = append(n.inboundConnections, conn)
	buffer := bufio.NewReader(conn)
	for {
		log.Println("reading data from inbound connection")
		data, err := buffer.ReadBytes('\n')
		if err == io.EOF {
			log.Printf("closing connection")
			conn.Close()
			return
		}
		if err != nil {
			log.Printf("connection error: %v", err)
			return
		}
		n.handleData(data[:len(data)-1])
	}
}

func (n *Node) handleListener(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("listener error: %v", err)
			break
		}
		go n.handleInboundConnection(conn)
	}
}

func (n *Node) handleOutboundConnection(conn net.Conn) {
	n.outboundConnections = append(n.outboundConnections, conn)
	//...
}

func (n *Node) initConnection(c Connection) {
	for {
		conn, err := c.Connect()
		if err != nil {
			log.Printf("failed to connect to remote server: %v", err)
			n.Errors <- err
			time.Sleep(time.Second * 10)
			continue
		}
		go n.handleOutboundConnection(conn)
		return
	}
}

func (n *Node) initConnections() {
	for _, c := range n.connections {
		go n.initConnection(c)
	}
}

func (n *Node) initListener(l Listener) {
	listener, err := l.Listen()
	if err != nil {
		log.Printf("failed to start listener: %v", err)
		n.Errors <- err
		return
	}
	go n.handleListener(listener)
}

func (n *Node) initListeners() {
	for _, l := range n.listeners {
		go n.initListener(l)
	}
}
