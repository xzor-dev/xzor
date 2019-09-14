package command

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/module/messenger"
	"github.com/xzor-dev/xzor/internal/xzor/command"
)

// CreateBoardName is the name of the CreateBoard command.
const CreateBoardName = "create-board"

// CreateBoard handles the creation of new boards.
type CreateBoard struct {
	Service *messenger.Service
}

// Execute uses the provided arguments to create a new board.
func (c *CreateBoard) Execute(args []interface{}) (*command.Response, error) {
	if len(args) < 1 {
		return nil, errors.New("invalid number of arguments")
	}

	title, ok := args[0].(string)
	if !ok {
		return nil, errors.New("could not convert argument to string")
	}

	board, err := c.Service.NewBoard(title)
	if err != nil {
		return nil, err
	}
	res := &command.Response{
		Value: board,
	}
	return res, nil
}

// Name returns the name of the command.
func (c *CreateBoard) Name() command.Name {
	return CreateBoardName
}
