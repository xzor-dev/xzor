package command

import (
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

// Execute runs the command with the provided parameters to create a new message.
func (cm *CreateMessage) Execute(params command.Params) (*command.Response, error) {
	threadHash, err := params.String("thread")
	if err != nil {
		return nil, err
	}
	body, err := params.String("body")
	if err != nil {
		return nil, err
	}
	message, err := cm.Service.NewMessage(messenger.ThreadHash(threadHash), body)
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
