package block

import (
	"encoding/json"

	"github.com/xzor-dev/xzor/internal/xzor/common"
)

// Block holds sequential data.
type Block struct {
	Data         interface{}
	Hash         Hash
	Index        Index
	PreviousHash Hash
	Timestamp    int64
}

// Hash is a unique string generated from a block.
type Hash string

// NewHash generates a unique hash for a block.
func NewHash(b *Block) (Hash, error) {
	byteData, err := json.Marshal(b.Data)
	if err != nil {
		return "", err
	}
	record := string(b.Index) + string(b.Timestamp) + string(byteData) + string(b.PreviousHash)
	hash, err := common.NewHash([]byte(record))
	if err != nil {
		return "", err
	}
	return Hash(hash), nil
}

// Index is used to order blocks within a chain.
type Index int

// Store handles storage operations for blocks.
type Store interface {
	Delete(Hash) error
	Read(Hash) (*Block, error)
	Write(*Block) error
}
