package command_test

import (
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/command"
)

func TestCommands(t *testing.T) {
	commander := &command.Commander{}
	command := &testCommand{
		name: "reverse",
		callback: func(data []byte) ([]byte, error) {
			res := make([]byte, len(data))
			for i, d := range data {
				j := len(data) - 1 - i
				res[j] = d
			}
			return res, nil
		},
	}
	commander.Register(command)

	testStr := "hello"
	expected := "olleh"
	res, err := commander.Execute("reverse", []byte(testStr))
	if err != nil {
		t.Fatalf("%v", err)
	}
	if string(res) != expected {
		t.Fatalf("unexpected result from command: wanted %s, got %s", expected, res)
	}
}

var _ command.Command = &testCommand{}

type testCommand struct {
	name     string
	callback func([]byte) ([]byte, error)
}

func (c *testCommand) Execute(data []byte) ([]byte, error) {
	return c.callback(data)
}

func (c *testCommand) Name() string {
	return c.name
}
