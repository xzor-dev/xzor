package command

import (
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
func (ct *CreateThread) Execute(params command.Params) (*command.Response, error) {
	boardHash, err := params.String("board")
	if err != nil {
		return nil, err
	}
	title, err := params.String("title")
	if err != nil {
		return nil, err
	}
	thread, err := ct.Service.NewThread(messenger.BoardHash(boardHash), title)
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
