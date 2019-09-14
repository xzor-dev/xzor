package module

import "github.com/xzor-dev/xzor/internal/xzor/command"

// Module defines what a system module requires.
type Module interface {
	Command(command.Name) (command.Command, error)
	Name() Name
}

// Name is a string name of a module.
type Name string
