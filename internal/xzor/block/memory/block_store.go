package memory

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/xzor/block"
)

var _ block.Store = &BlockStore{}

// BlockStore handles the storage of blocks within memory.
type BlockStore struct {
	blocks map[block.Hash]*block.Block
}

// Delete removes a block from the store.
func (s *BlockStore) Delete(hash block.Hash) error {
	if s.blocks != nil {
		delete(s.blocks, hash)
	}
	return nil
}

// Read attempts to get a block using its hash.
func (s *BlockStore) Read(hash block.Hash) (*block.Block, error) {
	if s.blocks == nil || s.blocks[hash] == nil {
		return nil, errors.New("invalid block hash")
	}
	return s.blocks[hash], nil
}

// Write adds or overwrites a block using its hash.
func (s *BlockStore) Write(b *block.Block) error {
	if s.blocks == nil {
		s.blocks = make(map[block.Hash]*block.Block)
	}
	if b.Hash == "" {
		return errors.New("block does not have a hash")
	}
	s.blocks[b.Hash] = b
	return nil
}
