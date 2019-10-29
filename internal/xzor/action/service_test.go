package action_test

import (
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/command"
	"github.com/xzor-dev/xzor/internal/xzor/module"
	"github.com/xzor-dev/xzor/internal/xzor/resource"
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
	providers := map[command.ProviderName]command.Provider{
		command.ProviderName(moduleName): m,
	}
	as := action.NewService(providers)
	a := &action.Action{
		Command:         commandName,
		CommandProvider: command.ProviderName(moduleName),
		Parameters: map[string]interface{}{
			"foo": "bar",
			"bar": "baz",
		},
	}
	_, err := as.ExecuteAction(a)
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
	a, err := action.New(command.ProviderName(mod.Name()), cmd.Name(), map[string]interface{}{"foo": "bar"})
	if err != nil {
		t.Fatalf("%v", err)
	}
	providers := map[command.ProviderName]command.Provider{
		"test-mod": mod,
	}
	s := action.NewService(providers)

	_, err = s.ExecuteAction(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	_, err = s.ExecuteAction(a)
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
var _ command.Provider = &testModule{}

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

func (m *testModule) Commands() map[command.Name]command.Command {
	return m.commands
}

func (m *testModule) Name() module.Name {
	return m.name
}

func (m *testModule) Resources() map[resource.Name]resource.Getter {
	return make(map[resource.Name]resource.Getter)
}
