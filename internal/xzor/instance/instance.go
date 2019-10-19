package instance

import (
	"errors"
	"log"

	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/module"
	"github.com/xzor-dev/xzor/internal/xzor/network"
)

// Instance ties all components together for a local running instance.
type Instance struct {
	messageHandlers []MessageHandler
	modules         map[module.Name]module.Module
	node            *network.Node
}

// New creates a new instance.
func New(modules []module.Module, node *network.Node) *Instance {
	moduleMap := make(map[module.Name]module.Module)
	for _, m := range modules {
		moduleMap[m.Name()] = m
	}
	return &Instance{
		messageHandlers: make([]MessageHandler, 0),
		modules:         moduleMap,
		node:            node,
	}
}

// ExecuteAction runs an action through the action service and propagates
// it through the network via the node.
func (i *Instance) ExecuteAction(a *action.Action) (*action.Response, error) {
	res, err := i.actionService.Execute(a)
	if err != nil {
		return nil, err
	}

	msg, err := network.NewMessage(a)
	if err != nil {
		return nil, err
	}

	err = i.node.Write(msg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Resource returns a single resource based on the supplied arguments.
func (i *Instance) Resource(moduleName module.Name, resourceName module.ResourceName, resourceID module.ResourceID) (module.Resource, error) {
	if i.modules == nil || i.modules[moduleName] == nil {
		return nil, errors.New("invalid module name provided")
	}

	getters := i.modules[moduleName].Resources()
	if getters[resourceName] == nil {
		return nil, errors.New("invalid resource name provided")
	}

	return getters[resourceName].Resource(resourceID)
}

// Start is used to start all necessary components within the instance.
func (i *Instance) Start() error {
	if i.node == nil {
		return errors.New("no node has been provided to the instance")
	}

	go i.readNodeMessages()

	return nil
}

// readNodeMessages reads incoming messages from the instance's node.
func (i *Instance) readNodeMessages() {
	for {
		data, err := i.node.Read()
		if err != nil {
			log.Printf("got error from node: %v", err)
		}

		action := &action.Action{}
		_, err = data.Decode(action)
		if err != nil {
			log.Printf("%v", err)
		}

		res, err := i.ExecuteAction(action)
		if err != nil {
			log.Printf("%v", err)
		} else {
			log.Printf("executed action %v", res.Action)
		}
	}
}

// MessageHandler is used to handle all incoming messages to the instance.
type MessageHandler interface {
	HandleMessage(network.EncodedMessage) error
}
