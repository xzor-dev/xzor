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
	return ""
}

var _ Connection = &MockConnection{}

// MockConnection implements Connection.
type MockConnection struct {
	Conn net.Conn
}

// Connect returns the net.Conn instance provided to the connection.
func (c *MockConnection) Connect() (net.Conn, error) {
	return c.Conn, nil
}

var _ DataHandler = &MockDataHandler{}

// MockDataHandler implements DataHandler using a custom handler function.
type MockDataHandler struct {
	Handler func([]byte) error
}

// HandleData passes the data to the handler function.
func (h *MockDataHandler) HandleData(data []byte) error {
	if h.Handler == nil {
		return errors.New("no handler function provided")
	}
	return h.Handler(data)
}

var _ Listener = &MockListener{}
var _ net.Listener = &MockListener{}

// MockListener implements both Listener and net.Listener.
type MockListener struct {
	Conn net.Conn

	connChan chan net.Conn
}

// Accept returns the connection provided to the listener.
// The connection is only returned once to prevent memory overflow.
func (l *MockListener) Accept() (net.Conn, error) {
	if l.connChan == nil {
		l.connChan = make(chan net.Conn)
		go func() {
			l.connChan <- l.Conn
		}()
	}
	return <-l.connChan, nil
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
