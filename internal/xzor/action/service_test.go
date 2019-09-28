package action_test

import (
	"errors"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/command"
	"github.com/xzor-dev/xzor/internal/xzor/module"
)

func TestIncomingAction(t *testing.T) {
	moduleName := module.Name("test-module")
	commandName := command.Name("test-command")

	c := &testCommand{
		name: commandName,
	}
	m := &testModule{
		name: moduleName,
		commands: map[command.Name]command.Command{
			commandName: c,
		},
	}
	as := &action.Service{
		Modules: map[module.Name]module.Module{
			moduleName: m,
		},
	}
	a := &action.Action{
		Arguments: []interface{}{"foo", "bar"},
		Command:   commandName,
		Module:    moduleName,
	}
	_, err := as.Execute(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(c.lastArgs) != 2 {
		t.Fatalf("expected 2 arguments from command, got %d", len(c.lastArgs))
	}
}

var _ command.Command = &testCommand{}

type testCommand struct {
	lastArgs []interface{}
	name     command.Name
}

func (c *testCommand) Execute(args []interface{}) (*command.Response, error) {
	c.lastArgs = args
	return &command.Response{
		Value: args,
	}, nil
}

func (c *testCommand) Name() command.Name {
	return c.name
}

var _ module.Module = &testModule{}

type testModule struct {
	commands map[command.Name]command.Command
	name     module.Name
}

func (m *testModule) Command(name command.Name) (command.Command, error) {
	if m.commands == nil || m.commands[name] == nil {
		return nil, errors.New("invalid command name")
	}
	return m.commands[name], nil
}

func (m *testModule) Name() module.Name {
	return m.name
}
