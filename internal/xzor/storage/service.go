package storage

import "errors"

// Service facilitates the creation and management of stored data.
type Service struct {
	ChainStore ChainStore
}

// NewChain creates new chain with a genesis block.
func (s *Service) NewChain() (*Chain, error) {
	hash, err := NewChainHash()
	if err != nil {
		return nil, err
	}
	c := &Chain{
		Blocks: make([]*Block, 0),
		Hash:   hash,
	}
	b := c.NewBlock(nil)
	err = c.AddBlock(b)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// DeleteChain deletes a chain from the chain store.
func (s *Service) DeleteChain(hash ChainHash) error {
	if s.ChainStore == nil {
		return errors.New("no ChainStore provided to the storage service")
	}
	return s.ChainStore.Delete(hash)
}

// ReadChain reads a chain from the chain store using its hash.
func (s *Service) ReadChain(hash ChainHash) (*Chain, error) {
	if s.ChainStore == nil {
		return nil, errors.New("no ChainStore provided to the storage service")
	}
	return s.ChainStore.Read(hash)
}

// WriteChain writes a chain to the chain store.
func (s *Service) WriteChain(c *Chain) error {
	if s.ChainStore == nil {
		return errors.New("no ChainStore provided to the storage service")
	}
	return s.ChainStore.Write(c)
}
