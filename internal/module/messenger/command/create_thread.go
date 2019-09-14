package command

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/module/messenger"
	"github.com/xzor-dev/xzor/internal/xzor/command"
)

// CreateThreadName is the name of the CreateThread command.
const CreateThreadName = "create-thread"

var _ command.Command = &CreateThread{}

// CreateThread handles the creation of new threads.
type CreateThread struct {
	Service *messenger.Service
}

// Execute the command using the provided arguments to create a new thread.
func (ct *CreateThread) Execute(args []interface{}) (*command.Response, error) {
	if ct.Service == nil {
		return nil, errors.New("no messenger service provided to command")
	}
	if len(args) < 2 {
		return nil, errors.New("expected at least two arguments")
	}

	boardHash, ok := args[0].(messenger.BoardHash)
	if !ok {
		return nil, errors.New("couldn't convert argument to BoardHash")
	}
	title, ok := args[1].(string)
	if !ok {
		return nil, errors.New("couldn't convert argument to string")
	}

	thread, err := ct.Service.NewThread(boardHash, title)
	if err != nil {
		return nil, err
	}

	return &command.Response{
		Value: thread,
	}, nil
}

// Name returns the name of the command.
func (ct *CreateThread) Name() command.Name {
	return CreateThreadName
}
