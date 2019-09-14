package command

import "errors"

// Command is used for executing single commands.
type Command interface {
	Execute([]interface{}) (*Response, error)
	Name() Name
}

// Commander groups commands and executes them.
type Commander struct {
	commands map[Name]Command
}

// Execute executes a command by its name along with the provided arguments.
func (c *Commander) Execute(name Name, args []interface{}) (*Response, error) {
	if c.commands == nil || c.commands[name] == nil {
		return nil, errors.New("invalid command name")
	}
	return c.commands[name].Execute(args)
}

// Command gets a command by its name.
func (c *Commander) Command(name Name) (Command, error) {
	if c.commands == nil || c.commands[name] == nil {
		return nil, errors.New("invalid command name")
	}
	return c.commands[name], nil
}

// Register adds or replaces a command by its name.
func (c *Commander) Register(cmd Command) {
	if c.commands == nil {
		c.commands = make(map[Name]Command)
	}
	c.commands[cmd.Name()] = cmd
}

// NewCommander creates a new Commander instance and registers the provided commands.
func NewCommander(commands []Command) *Commander {
	commander := &Commander{}
	for _, cmd := range commands {
		commander.Register(cmd)
	}
	return commander
}

// Name is a string name of a command.
type Name string

// Response is populated and returned from executed commands.
type Response struct {
	Value interface{}
}
