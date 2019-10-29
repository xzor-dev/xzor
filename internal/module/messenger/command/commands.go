package command

import (
	"github.com/xzor-dev/xzor/internal/module/messenger"
	"github.com/xzor-dev/xzor/internal/xzor/command"
)

// Commands creates a slice of available commands for the messenger module.
func Commands(s *messenger.Service) []command.Command {
	return []command.Command{
		&CreateBoard{
			Service: s,
		},
		&CreateMessage{
			Service: s,
		},
		&CreateThread{
			Service: s,
		},
	}
}
