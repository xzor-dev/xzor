package messenger

import (
	"github.com/xzor-dev/xzor/internal/xzor/common"
)

// Thread holds a collection of messages.
type Thread struct {
	Hash     ThreadHash
	Messages []MessageHash
	Title    string
}

// HasMessage checks if the thread has a message.
func (t *Thread) HasMessage(hash MessageHash) bool {
	if t.Messages == nil {
		return false
	}
	for _, h := range t.Messages {
		if h == hash {
			return true
		}
	}
	return false
}

// NewMessage creates a new message in the thread.
func (t *Thread) NewMessage(body string) (*Message, error) {
	if t.Messages == nil {
		t.Messages = make([]MessageHash, 0)
	}
	hash, err := NewMessageHash()
	if err != nil {
		return nil, err
	}
	t.Messages = append(t.Messages, hash)
	return &Message{
		Body: body,
		Hash: hash,
	}, nil
}

// ThreadHash is a unique string hash assigned to threads.
type ThreadHash string

// NewThreadHash generates a new thread hash.
func NewThreadHash() (ThreadHash, error) {
	rb, err := common.NewRandomBytes(32)
	if err != nil {
		return "", err
	}
	hash, err := common.NewHash(rb)
	if err != nil {
		return "", err
	}
	return ThreadHash(hash), nil
}
