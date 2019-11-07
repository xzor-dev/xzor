package simulator

import (
	"net"

	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/network"
)

// Network holds the simulator network of nodes.
type Network struct {
	Nodes []*network.Node
}

// NewNetwork generates a new network map using the provided configuration and builder.
func NewNetwork(nodes []*network.Node) (*Network, error) {
	return &Network{
		Nodes: nodes,
	}, nil
}

// Node returns a node from the network at the specified index.
func (n *Network) Node(index int) (*network.Node, error) {
	if n.Nodes[index] == nil {
		return nil, ErrInvalidNodeIndex
	}
	return n.Nodes[index], nil
}

// Push sends an action to the root node.
func (n *Network) Push(a *action.Action) error {
	node, err := n.Node(0)
	if err != nil {
		return err
	}
	return node.Write(a)
}

// NetworkBuilder is used to build network configurations.
type NetworkBuilder interface {
	Build() ([]*network.Node, error)
}

// NetworkLoopBuilder is used to build a network map where
// nodes eventually loop back to the first node.
type NetworkLoopBuilder struct {
	MaxConnections int
	TotalNodes     int
}

// NewNetworkLoopBuilder creates a new NetworkLoopBuilder using
// the supplied arguments.
func NewNetworkLoopBuilder(totalNodes int, maxConnections int) *NetworkLoopBuilder {
	return &NetworkLoopBuilder{
		MaxConnections: maxConnections,
		TotalNodes:     totalNodes,
	}
}

// Build generates a looped network map of nodes.
func (b *NetworkLoopBuilder) Build() ([]*network.Node, error) {
	nodes := make([]*network.Node, b.TotalNodes)
	listeners := make([]*network.MockListener, b.TotalNodes)
	for i := 0; i < b.TotalNodes; i++ {
		nodes[i] = network.NewNode()
		listeners[i] = &network.MockListener{}
	}

	for i := 0; i < b.TotalNodes; i++ {
		nodeA := nodes[i]

		for j := 0; j < b.MaxConnections; j++ {
			next := i + j + 1
			if next >= b.TotalNodes {
				diff := next - b.TotalNodes
				if diff == i {
					diff++
				}
				next = diff
			}

			listenerB := listeners[next]
			pipeA, pipeB := net.Pipe()

			nodeA.AddConnection(pipeA)
			listenerB.AddConnection(pipeB)
		}
	}
	for i := 0; i < b.TotalNodes; i++ {
		nodes[i].AddListener(listeners[i])
	}

	return nodes, nil
}

// NetworkWebBuilder builds a web-like network map.
type NetworkWebBuilder struct {
	MaxConnections int
	TotalNodes     int
}

// NewNetworkWebBuilder creates a new NetworkWebBuilder using
// the supplied arguments.
func NewNetworkWebBuilder(totalNodes int, maxConnections int) *NetworkWebBuilder {
	return &NetworkWebBuilder{
		MaxConnections: maxConnections,
		TotalNodes:     totalNodes,
	}
}

// Build generates a slice of nodes arranged in a web configuration.
func (b *NetworkWebBuilder) Build() ([]*network.Node, error) {
	nodes := make([]*network.Node, b.TotalNodes)
	listeners := make([]*network.MockListener, b.TotalNodes)

	for i := 0; i < b.TotalNodes; i++ {
		nodes[i] = network.NewNode()
		listeners[i] = &network.MockListener{}
	}

	for i := 0; i < b.TotalNodes; i++ {
		nodeA := nodes[i]
		listenerA := listeners[i]

		for j := 0; j < b.MaxConnections; j++ {
			k := i*b.MaxConnections + j + 1
			if len(nodes) < k+1 {
				break
			}

			nodeB := nodes[k]
			listenerB := listeners[k]

			pipeA1, pipeA2 := net.Pipe()
			pipeB1, pipeB2 := net.Pipe()

			nodeA.AddConnection(pipeB2)
			listenerA.AddConnection(pipeA1)

			nodeB.AddConnection(pipeA2)
			listenerB.AddConnection(pipeB1)
		}

		nodeA.AddListener(listenerA)
	}

	return nodes, nil
}
