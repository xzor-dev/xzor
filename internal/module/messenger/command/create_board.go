package command

import (
	"github.com/xzor-dev/xzor/internal/module/messenger"
	"github.com/xzor-dev/xzor/internal/xzor/command"
)

// CreateBoardName is the name of the CreateBoard command.
const CreateBoardName = "create-board"

var _ command.Command = &CreateBoard{}

// CreateBoard handles the creation of new boards.
type CreateBoard struct {
	Service *messenger.Service
}

// Execute uses the provided parameters to create a new board.
func (c *CreateBoard) Execute(params command.Params) (*command.Response, error) {
	title, err := params.String("title")
	if err != nil {
		return nil, err
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
