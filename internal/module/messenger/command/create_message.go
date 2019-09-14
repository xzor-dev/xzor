package command

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/module/messenger"
	"github.com/xzor-dev/xzor/internal/xzor/command"
)

// CreateMessageName is the name of the CreateMessage command.
const CreateMessageName = "create-message"

var _ command.Command = &CreateMessage{}

// CreateMessage is used to create messages within threads.
type CreateMessage struct {
	Service *messenger.Service
}

// Execute runs the command with the provided arguments to create a new message.
func (cm *CreateMessage) Execute(args []interface{}) (*command.Response, error) {
	if len(args) < 2 {
		return nil, errors.New("command requires at 2 arguments")
	}

	threadHash, ok := args[0].(messenger.ThreadHash)
	if !ok {
		return nil, errors.New("could not convert argument to ThreadHash")
	}
	body, ok := args[1].(string)
	if !ok {
		return nil, errors.New("could not convert argument to string")
	}

	message, err := cm.Service.NewMessage(threadHash, body)
	if err != nil {
		return nil, err
	}
	return &command.Response{
		Value: message,
	}, nil
}

// Name returns the name of the command.
func (cm *CreateMessage) Name() command.Name {
	return CreateMessageName
}
