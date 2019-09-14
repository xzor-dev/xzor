package action

import (
	"github.com/xzor-dev/xzor/internal/xzor/command"
	"github.com/xzor-dev/xzor/internal/xzor/module"
)

// Action contains a single command along with any arguments.
type Action struct {
	Arguments []interface{}
	Command   command.Name
	Module    module.Name
}

// Decoder converts a byte slice to an action.
type Decoder interface {
	DecodeAction([]byte) (*Action, error)
}

// Encoder converts an action into a byte slice.
type Encoder interface {
	EncodeAction(*Action) ([]byte, error)
}

// EncodeDecoder combines the Decoder and Encoder interfaces.
type EncodeDecoder interface {
	Decoder
	Encoder
}

// Response is populated and returned from executing actions.
type Response struct {
	Action *Action
	Value  interface{}
}
