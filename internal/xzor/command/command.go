package command

import "errors"

type Command interface {
	Execute([]byte) ([]byte, error)
	Name() string
}

type Commander struct {
	commands map[string]Command
}

func (c *Commander) Execute(name string, data []byte) ([]byte, error) {
	if c.commands == nil || c.commands[name] == nil {
		return nil, errors.New("invalid command name")
	}
	return c.commands[name].Execute(data)
}

func (c *Commander) Register(cmd Command) {
	if c.commands == nil {
		c.commands = make(map[string]Command)
	}
	c.commands[cmd.Name()] = cmd
}
