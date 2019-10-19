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
	as := action.NewService([]module.Module{m})
	a := &action.Action{
		Command: commandName,
		Module:  moduleName,
		Parameters: map[string]interface{}{
			"foo": "bar",
			"bar": "baz",
		},
	}
	_, err := as.Execute(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(c.lastParams) != 2 {
		t.Fatalf("expected 2 arguments from command, got %d", len(c.lastParams))
	}
}

func TestDuplicateActionIgnore(t *testing.T) {
	cmd := newTestCommand("test-cmd")
	mod := newTestModule("test-mod", []command.Command{cmd})
	a, err := action.New(mod.Name(), cmd.Name(), map[string]interface{}{"foo": "bar"})
	if err != nil {
		t.Fatalf("%v", err)
	}
	s := action.NewService([]module.Module{mod})

	_, err = s.Execute(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	_, err = s.Execute(a)
	if err != action.ErrDuplicateAction {
		t.Fatalf("expected a duplicate action error, got %v", err)
	}
}

var _ command.Command = &testCommand{}

type testCommand struct {
	lastParams map[string]interface{}
	name       command.Name
}

func newTestCommand(name command.Name) *testCommand {
	return &testCommand{
		name: name,
	}
}

func (c *testCommand) Execute(params map[string]interface{}) (*command.Response, error) {
	c.lastParams = params
	return &command.Response{
		Value: params,
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

func newTestModule(name module.Name, commands []command.Command) *testModule {
	mod := &testModule{
		commands: make(map[command.Name]command.Command),
		name:     name,
	}
	for _, cmd := range commands {
		mod.commands[cmd.Name()] = cmd
	}
	return mod
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

func (m *testModule) Resources() map[module.ResourceName]module.ResourceGetter {
	return make(map[module.ResourceName]module.ResourceGetter)
}
