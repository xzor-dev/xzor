package svc

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	// ModuleNameDelimiter is used to separate a module's name
	// from the data to be processed by that module.
	ModuleNameDelimiter = ":"
)

// ModuleName is a unique identifier string for modules.
type ModuleName string

// Module defines unique handlers of incoming data.
type Module interface {
	// Name returns the module's unique name.
	Name() ModuleName

	// Process handles processing of data related to the module.
	Process([]byte) error
}

// ModuleRouter takes incoming messages and forwards the data to
// the appropriate module.
type ModuleRouter struct {
	Modules map[ModuleName]Module
}

// Process parses the provided data and passes it to the appropriate module.
func (r *ModuleRouter) Process(msg []byte) error {
	if r.Modules == nil {
		return errors.New("no modules have been registered with the router")
	}

	parts := bytes.Split(msg, []byte(ModuleNameDelimiter))
	if len(parts) != 2 {
		return errors.New("module name could not be parsed from message data")
	}

	moduleName := ModuleName(parts[0])
	data := parts[1]
	if r.Modules[moduleName] == nil {
		return fmt.Errorf("invalid module name: %s", moduleName)
	}
	return r.Modules[moduleName].Process(data)
}

// Register adds a module to the router.
func (r *ModuleRouter) Register(m Module) {
	if r.Modules == nil {
		r.Modules = make(map[ModuleName]Module)
	}
	r.Modules[m.Name()] = m
}
