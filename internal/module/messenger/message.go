package messenger

import "github.com/xzor-dev/xzor/internal/xzor/common"

// Message holds a single message.
type Message struct {
	Body string
	Hash MessageHash
}

// MessageHash is a unique hash string assigned to messages.
type MessageHash string

// NewMessageHash creates a new hash for a message.
func NewMessageHash() (MessageHash, error) {
	rb, err := common.NewRandomBytes(32)
	if err != nil {
		return "", err
	}
	hash, err := common.NewHash(rb)
	if err != nil {
		return "", err
	}
	return MessageHash(hash), nil
}
