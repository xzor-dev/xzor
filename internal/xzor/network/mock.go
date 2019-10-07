package network

import (
	"errors"
	"net"
)

var _ net.Addr = &MockAddr{}

// MockAddr implements net.Addr.
type MockAddr struct{}

// Network returns the type of network.
func (a *MockAddr) Network() string {
	return "mock"
}

// String returns the address as a string.
func (a *MockAddr) String() string {
	return "mock"
}

var _ net.Listener = &MockListener{}

// MockListener implements both Listener and net.Listener.
type MockListener struct {
	Connections []net.Conn
}

// Accept returns the next connection in the Connections slice.
func (l *MockListener) Accept() (net.Conn, error) {
	if len(l.Connections) == 0 {
		return nil, errors.New("no more connections")
	}
	next, conns := l.Connections[0], l.Connections[1:]
	l.Connections = conns
	return next, nil
}

// AddConnection adds a new net.Conn to the connection stack.
// Connections are returned in the order they were added when calling Accept().
func (l *MockListener) AddConnection(conn net.Conn) {
	l.Connections = append(l.Connections, conn)
}

// Addr returns the net.Addr for the listener.
func (l *MockListener) Addr() net.Addr {
	return &MockAddr{}
}

// Close closes the listener.
func (l *MockListener) Close() error {
	return nil
}

// Listen returns the MockListener as a net.Listener instance.
func (l *MockListener) Listen() (net.Listener, error) {
	return l, nil
}
