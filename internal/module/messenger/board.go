package messenger

import (
	"github.com/xzor-dev/xzor/internal/xzor/common"
)

// Board holds threads.
type Board struct {
	Hash    BoardHash
	Threads []ThreadHash
	Title   string
}

// HasThread checks if the board has a thread.
func (b *Board) HasThread(hash ThreadHash) bool {
	if b.Threads == nil {
		return false
	}
	for _, h := range b.Threads {
		if h == hash {
			return true
		}
	}
	return false
}

// NewThread creates a new thread on the board with the provided title.
func (b *Board) NewThread(title string) (*Thread, error) {
	if b.Threads == nil {
		b.Threads = make([]ThreadHash, 0)
	}
	hash, err := NewThreadHash()
	if err != nil {
		return nil, err
	}
	thread := &Thread{
		Hash:  hash,
		Title: title,
	}
	b.Threads = append(b.Threads, hash)
	return thread, nil
}

// BoardHash is a unique string assigned to newly created boards.
type BoardHash string

// NewBoardHash generates hashes for new boards.
func NewBoardHash() (BoardHash, error) {
	rb, err := common.NewRandomBytes(32)
	if err != nil {
		return "", err
	}
	hash, err := common.NewHash(rb)
	if err != nil {
		return "", err
	}
	return BoardHash(hash), nil
}
