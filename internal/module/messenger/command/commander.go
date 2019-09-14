package command

import (
	"github.com/xzor-dev/xzor/internal/module/messenger"
	"github.com/xzor-dev/xzor/internal/xzor/command"
)

// NewCommander creates a new Commander instance populated with messenger commands.
func NewCommander(s *messenger.Service) *command.Commander {
	return command.NewCommander([]command.Command{
		&CreateBoard{
			Service: s,
		},
		&CreateMessage{
			Service: s,
		},
		&CreateThread{
			Service: s,
		},
	})
}
