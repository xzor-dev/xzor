package network

import (
	"encoding/json"
	"time"

	"github.com/xzor-dev/xzor/internal/xzor/common"
)

// Message holds data that is transported to and from nodes.
type Message struct {
	Data      interface{}
	Hash      MessageHash
	Timestamp int64
}

// NewMessage creates a new message instance with the supplied data.
func NewMessage(data interface{}) (*Message, error) {
	hash, err := NewMessageHash()
	if err != nil {
		return nil, err
	}
	m := &Message{
		Data:      data,
		Hash:      hash,
		Timestamp: time.Now().Unix(),
	}
	return m, nil
}

// Encode converts the message to a JSON byte slice as an EncodedMessage.
func (m *Message) Encode() (EncodedMessage, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return EncodedMessage(b), nil
}

// MessageHash is a unique string hash of a message.
type MessageHash string

// NewMessageHash generates a new hash for a message.
func NewMessageHash() (MessageHash, error) {
	bytes, err := common.NewRandomBytes(32)
	if err != nil {
		return "", err
	}
	hash, err := common.NewHash(bytes)
	if err != nil {
		return "", err
	}
	return MessageHash(hash), nil
}

// EncodedMessage is a JSON encoded message that is sent to and recieved by network nodes.
type EncodedMessage []byte

// Decode converts the encoded message back into a message struct using
// the supplied data interface as the decoded message's data property.
func (d EncodedMessage) Decode(data interface{}) (*Message, error) {
	msg := &Message{
		Data: data,
	}
	err := json.Unmarshal(d, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
