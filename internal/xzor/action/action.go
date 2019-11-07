package action

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/xzor-dev/xzor/internal/xzor/command"
	"github.com/xzor-dev/xzor/internal/xzor/common"
)

// Action contains a single command along with any arguments.
type Action struct {
	Command         command.Name
	CommandProvider command.ProviderName
	Hash            Hash
	Parameters      map[string]interface{}
	Timestamp       int64
}

// New generates a new action.
func New(providerName command.ProviderName, cmd command.Name, params map[string]interface{}) (*Action, error) {
	t := time.Now()
	hash, err := NewHash(providerName, cmd, params, t)
	if err != nil {
		return nil, err
	}
	return &Action{
		Command:         cmd,
		CommandProvider: providerName,
		Hash:            hash,
		Parameters:      params,
		Timestamp:       t.Unix(),
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
func NewHash(providerName command.ProviderName, cmd command.Name, params map[string]interface{}, t time.Time) (Hash, error) {
	pb, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	hs := fmt.Sprintf("%s-%s-%s-%d", providerName, cmd, pb, t.Unix())
	hb := []byte(hs)
	hash, err := common.NewHash(hb)
	if err != nil {
		return "", err
	}
	log.Printf("made hash '%s' from string '%s'", hash, hs)
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

// Executor is used to execute actions and return an action response.
type Executor interface {
	ExecuteAction(*Action) (*Response, error)
}

// Response is populated and returned from executing actions.
type Response struct {
	Action *Action
	Value  interface{}
}
