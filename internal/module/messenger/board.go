package messenger

import (
	"github.com/xzor-dev/xzor/internal/xzor/common"
	"github.com/xzor-dev/xzor/internal/xzor/module"
)

// BoardResourceName is the name of the Board resource.
const BoardResourceName = "board"

var _ module.Resource = &Board{}

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

// ResourceName returns the resource name of the board.
func (b *Board) ResourceName() module.ResourceName {
	return BoardResourceName
}

// ResourceID returns the board's hash as a ResourceID.
func (b *Board) ResourceID() module.ResourceID {
	return module.ResourceID(b.Hash)
}

// BoardHash is a unique string assigned to newly created boards.
type BoardHash string

// NewBoardHash generates hashes for new boards.
func NewBoardHash(title string) (BoardHash, error) {
	b := []byte(title)
	hash, err := common.NewHash(b)
	if err != nil {
		return "", err
	}
	return BoardHash(hash), nil
}

var _ module.ResourceGetter = &BoardResourceGetter{}

// BoardResourceGetter handles the retrieval of individual boards.
type BoardResourceGetter struct {
	service *Service
}

// NewBoardResourceGetter creates a new instance of BoardResourceGetter with the supplied service.
func NewBoardResourceGetter(s *Service) *BoardResourceGetter {
	return &BoardResourceGetter{
		service: s,
	}
}

// Resource returns a board with the hash of the id supplied.
func (g *BoardResourceGetter) Resource(id module.ResourceID) (module.Resource, error) {
	return g.service.Board(BoardHash(id))
}
