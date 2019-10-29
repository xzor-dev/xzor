package command_test

import (
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/command"
)

func TestCommands(t *testing.T) {
	commander := &command.Commander{}
	c := &testCommand{
		name: "reverse",
		callback: func(params map[string]interface{}) ([]byte, error) {
			value := params["value"].(string)
			data := []byte(value)
			res := make([]byte, len(data))
			for i, d := range data {
				j := len(data) - 1 - i
				res[j] = d
			}
			return res, nil
		},
	}
	commander.Register(c)

	testStr := "hello"
	expected := "olleh"
	params := map[string]interface{}{
		"value": testStr,
	}
	res, err := commander.Execute("reverse", params)
	if err != nil {
		t.Fatalf("%v", err)
	}
	resValue, ok := res.Value.([]byte)
	if !ok {
		t.Fatal("couldn't convert response value to byte slice")
	}
	if string(resValue) != expected {
		t.Fatalf("unexpected result from command: wanted %s, got %s", expected, resValue)
	}
}

var _ command.Command = &testCommand{}

type testCommand struct {
	name     command.Name
	callback func(map[string]interface{}) ([]byte, error)
}

func (c *testCommand) Execute(params command.Params) (*command.Response, error) {
	res, err := c.callback(params)
	if err != nil {
		return nil, err
	}
	return &command.Response{
		Value: res,
	}, nil
}

func (c *testCommand) Name() command.Name {
	return c.name
}
