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
	Service *Service
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