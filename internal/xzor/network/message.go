package network

import (
	"time"

	"github.com/xzor-dev/xzor/internal/xzor/common"
)

// Message holds data that is transported to and from nodes.
type Message struct {
	Data      interface{} `json:"data"`
	Hash      MessageHash `json:"hash"`
	Timestamp int64       `json:"timestamp"`
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
