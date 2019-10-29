package instance

import (
	"errors"
	"log"

	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/network"
	"github.com/xzor-dev/xzor/internal/xzor/resource"
)

// Instance ties all components together for a local running instance.
type Instance struct {
	actionExecutor  action.Executor
	node            *network.Node
	resourceService *resource.Service
}

// New creates a new instance.
func New(actionExecutor action.Executor, node *network.Node, resourceService *resource.Service) *Instance {
	return &Instance{
		actionExecutor:  actionExecutor,
		node:            node,
		resourceService: resourceService,
	}
}

// ExecuteAction runs an action through the action executor and propagates
// it through the network via the node.
func (i *Instance) ExecuteAction(a *action.Action) (*action.Response, error) {
	res, err := i.actionExecutor.ExecuteAction(a)
	if err != nil {
		return nil, err
	}

	err = i.node.Write(a)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Resource returns a single resource based on the supplied arguments.
func (i *Instance) Resource(providerName resource.ProviderName, resourceName resource.Name, resourceID resource.ID) (resource.Resource, error) {
	return i.resourceService.Resource(providerName, resourceName, resourceID)
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
		action, err := i.node.Read()
		if err != nil {
			log.Printf("got error from node: %v", err)
		}

		res, err := i.ExecuteAction(action)
		if err != nil {
			log.Printf("%v", err)
		} else {
			log.Printf("executed action %v", res.Action)
		}
	}
}
