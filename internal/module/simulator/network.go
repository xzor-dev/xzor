package simulator

import (
	"net"

	"github.com/xzor-dev/xzor/internal/xzor/network"
)

// Network holds the simulator network of nodes.
type Network struct {
	Nodes []*network.Node
}

// NewNetwork generates a new network map using the provided configuration and builder.
func NewNetwork(config *Config, builder NetworkBuilder) (*Network, error) {
	nodes, err := builder.Build(config)
	if err != nil {
		return nil, err
	}
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

// NetworkBuilder is used to build network configurations.
type NetworkBuilder interface {
	Build(*Config) ([]*network.Node, error)
}

// NetworkLoopBuilder is used to build a network map where
// nodes eventually loop back to the first node.
type NetworkLoopBuilder struct{}

// Build generates a looped network map of nodes.
func (b *NetworkLoopBuilder) Build(c *Config) ([]*network.Node, error) {
	if c.ConnectionsPerNode >= c.TotalNodes {
		return nil, ErrNetworkSizeTooSmall
	}
	nodes := make([]*network.Node, c.TotalNodes)
	listeners := make([]*network.MockListener, c.TotalNodes)
	for i := 0; i < c.TotalNodes; i++ {
		nodes[i] = network.NewNode()
		listeners[i] = &network.MockListener{}
	}

	for i := 0; i < c.TotalNodes; i++ {
		nodeA := nodes[i]

		for j := 0; j < c.ConnectionsPerNode; j++ {
			next := i + j + 1
			if next >= c.TotalNodes {
				diff := next - c.TotalNodes
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
	for i := 0; i < c.TotalNodes; i++ {
		nodes[i].AddListener(listeners[i])
	}

	return nodes, nil
}

// NetworkWebBuilder builds a web-like network map.
type NetworkWebBuilder struct{}

// Build generates a slice of nodes arranged in a web configuration.
func (b *NetworkWebBuilder) Build(c *Config) ([]*network.Node, error) {
	nodes := make([]*network.Node, c.TotalNodes)
	listeners := make([]*network.MockListener, c.TotalNodes)

	for i := 0; i < c.TotalNodes; i++ {
		nodes[i] = network.NewNode()
		listeners[i] = &network.MockListener{}
	}

	for i := 0; i < c.TotalNodes; i++ {
		nodeA := nodes[i]
		listenerA := listeners[i]

		for j := 0; j < c.ConnectionsPerNode; j++ {
			k := i*c.ConnectionsPerNode + j + 1
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
