package block

import "errors"

// Service facilitates the creation and management of stored data.
type Service struct {
	BlockStore Store
	ChainStore ChainStore
}

// NewBlock creates a new block for the provided chain and
// guarantees it as a valid next block on the chain.
func (s *Service) NewBlock(c *Chain, data []byte) (*Block, error) {
	for {
		b := c.NewBlock(nil)
		err := c.AddBlock(b)
		if err == ErrInvalidPrevHash {
			continue
		} else if err != nil {
			return nil, err
		}
		return b, nil
	}
}

func (s *Service) DeleteBlock(hash Hash) error {
	return s.BlockStore.Delete(hash)
}

func (s *Service) ReadBlock(hash Hash) (*Block, error) {
	return s.BlockStore.Read(hash)
}

func (s *Service) WriteBlock(b *Block) error {
	return s.BlockStore.Write(b)
}

// NewBranch creates a branch from a chain's block.
func (s *Service) NewBranch(fromChain *Chain, fromBlock *Block) (*Branch, error) {
	c2, err := s.NewChain()
	if err != nil {
		return nil, err
	}

	branch, err := fromChain.NewBranch(fromBlock, c2)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

// DeleteChain deletes a chain from the chain store.
func (s *Service) DeleteChain(hash ChainHash) error {
	if s.ChainStore == nil {
		return errors.New("no ChainStore provided to the storage service")
	}
	return s.ChainStore.Delete(hash)
}

// NewChain creates new chain with a genesis block.
func (s *Service) NewChain() (*Chain, error) {
	hash, err := NewChainHash()
	if err != nil {
		return nil, err
	}
	c := &Chain{
		Blocks:   make(map[Hash]Index),
		Branches: make(map[BranchHash]*Branch),
		Hash:     hash,
	}
	b := c.NewBlock(nil)
	err = c.AddBlock(b)
	if err != nil {
		return nil, err
	}
	return c, nil
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
