package memory

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/xzor/block"
)

var _ block.ChainStore = &ChainStore{}

// ChainStore implements block.ChainStore to store chain data in memory.
type ChainStore struct {
	chains map[block.ChainHash]*block.Chain
}

// Delete removes a chain from storage using its hash.
func (s *ChainStore) Delete(hash block.ChainHash) error {
	if s.chains != nil {
		delete(s.chains, hash)
	}
	return nil
}

// Read attempts to get a chain from memory using its hash.
func (s *ChainStore) Read(hash block.ChainHash) (*block.Chain, error) {
	if s.chains == nil || s.chains[hash] == nil {
		return nil, errors.New("invalid chain hash")
	}
	return s.chains[hash], nil
}

// Write adds or replaces a chain in memory.
func (s *ChainStore) Write(c *block.Chain) error {
	if s.chains == nil {
		s.chains = make(map[block.ChainHash]*block.Chain)
	}
	s.chains[c.Hash] = c
	return nil
}
