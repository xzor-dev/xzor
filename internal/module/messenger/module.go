package messenger

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/xzor/command"
	"github.com/xzor-dev/xzor/internal/xzor/module"
)

// ModuleName is the name of the messenger module.
const ModuleName = "messenger"

var _ module.Module = &Module{}

// Module implements module.Module.
type Module struct {
	Commander *command.Commander
	Service   *Service

	resourceGetters map[module.ResourceName]module.ResourceGetter
}

// NewModule creates a new messenger module instance.
func NewModule(service *Service, commander *command.Commander) *Module {
	return &Module{
		Commander: commander,
		Service:   service,

		resourceGetters: map[module.ResourceName]module.ResourceGetter{
			BoardResourceName: &BoardResourceGetter{service},
		},
	}
}

// Command gets a messenger command by its name.
func (m *Module) Command(name command.Name) (command.Command, error) {
	if m.Commander == nil {
		return nil, errors.New("no commander provided to the module")
	}
	return m.Commander.Command(name)
}

// Name returns the name of the messenger module.
func (m *Module) Name() module.Name {
	return ModuleName
}

// Resources returns a map of all resource getters for the module.
func (m *Module) Resources() map[module.ResourceName]module.ResourceGetter {
	return m.resourceGetters
}
