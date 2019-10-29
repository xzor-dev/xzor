package messenger

import (
	"github.com/xzor-dev/xzor/internal/xzor/command"
	"github.com/xzor-dev/xzor/internal/xzor/module"
	"github.com/xzor-dev/xzor/internal/xzor/resource"
)

// ModuleName is the name of the messenger module.
const ModuleName = "messenger"

var _ command.Provider = &Module{}
var _ module.Module = &Module{}
var _ resource.Provider = &Module{}

// Module implements module.Module.
type Module struct {
	commands        command.Map
	resourceGetters resource.GetterMap
	service         *Service
}

// NewModule creates a new messenger module instance.
func NewModule(service *Service, commands []command.Command) *Module {
	return &Module{
		commands: command.NewMap(commands),
		resourceGetters: resource.GetterMap{
			BoardResourceName: &BoardResourceGetter{service},
		},
		service: service,
	}
}

// Commands returns a map of supported commands.
func (m *Module) Commands() command.Map {
	return m.commands
}

// CommandProviderName returns the name of the module as a command provider name.
func (m *Module) CommandProviderName() command.ProviderName {
	return ModuleName
}

// Name returns the name of the messenger module.
func (m *Module) Name() module.Name {
	return ModuleName
}

// Resources returns a map of all resource getters for the module.
func (m *Module) Resources() resource.GetterMap {
	return m.resourceGetters
}

// ResourceProviderName returns the name of the module as a resource provider name
func (m *Module) ResourceProviderName() resource.ProviderName {
	return ModuleName
}
