package action

import (
	"encoding/json"
	"time"

	"github.com/xzor-dev/xzor/internal/xzor/command"
	"github.com/xzor-dev/xzor/internal/xzor/common"
	"github.com/xzor-dev/xzor/internal/xzor/module"
)

// Action contains a single command along with any arguments.
type Action struct {
	Command    command.Name
	Hash       Hash
	Module     module.Name
	Parameters map[string]interface{}
	Timestamp  int64
}

// New generates a new action.
func New(mod module.Name, cmd command.Name, params map[string]interface{}) (*Action, error) {
	t := time.Now()
	hash, err := NewHash(mod, cmd, params, t)
	if err != nil {
		return nil, err
	}
	return &Action{
		Command:    cmd,
		Hash:       hash,
		Module:     mod,
		Parameters: params,
		Timestamp:  t.Unix(),
	}, nil
}

// Encode converts the action into a JSON string as EncodedAction.
func (a *Action) Encode() (EncodedAction, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return EncodedAction(data), nil
}

// Hash is a unique hash of an action.
type Hash string

// NewHash generates a new hash used for an action.
func NewHash(mod module.Name, cmd command.Name, params map[string]interface{}, t time.Time) (Hash, error) {
	pb, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	hs := string(mod) + string(cmd) + string(pb) + string(t.Unix())
	hb := []byte(hs)
	hash, err := common.NewHash(hb)
	if err != nil {
		return "", err
	}
	return Hash(hash), nil
}

// EncodedAction is an action encoded as a JSON byte slice.
type EncodedAction []byte

// Decode converts an encoded action back into an action struct.
func (en EncodedAction) Decode() (*Action, error) {
	a := &Action{}
	err := json.Unmarshal(en, a)
	return a, err
}

// Handler is used to handle actions.
type Handler interface {
	HandleAction(*Action) error
}

/*
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
*/

// Response is populated and returned from executing actions.
type Response struct {
	Action *Action
	Value  interface{}
}
