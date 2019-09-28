package block

import "errors"

// Service facilitates the creation and management of blocks.
type Service struct {
	BlockStore Store
	ChainStore ChainStore

	blocks map[Hash]*Block
	chains map[ChainHash]*Chain
}

// Block reads a block from the block store using its hash.
func (s *Service) Block(hash Hash) (*Block, error) {
	if s.BlockStore == nil {
		return nil, errors.New("no BlockStore provided to the service")
	}
	return s.BlockStore.Read(hash)
}

// Chain reads a chain from the chain store using its hash.
func (s *Service) Chain(hash ChainHash) (*Chain, error) {
	if s.ChainStore == nil {
		return nil, errors.New("no ChainStore provided to the service")
	}
	return s.ChainStore.Read(hash)
}

// Commit writes any blocks or chains in memory to their respective storages.
func (s *Service) Commit() error {
	if s.blocks != nil {
		for _, b := range s.blocks {
			err := s.WriteBlock(b)
			if err != nil {
				return err
			}
		}
	}
	if s.chains != nil {
		for _, c := range s.chains {
			err := s.WriteChain(c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteBlock removes a block from memory and the block store.
func (s *Service) DeleteBlock(hash Hash) error {
	err := s.BlockStore.Delete(hash)
	if err != nil {
		return err
	}
	if s.blocks != nil {
		delete(s.blocks, hash)
	}
	return nil
}

// DeleteChain deletes a chain from the chain store along with
// all blocks within it.
func (s *Service) DeleteChain(hash ChainHash) error {
	if s.ChainStore == nil {
		return errors.New("no ChainStore provided to the storage service")
	}
	c, err := s.ChainStore.Read(hash)
	if err != nil {
		return err
	}
	for blockHash := range c.Blocks {
		err := s.DeleteBlock(blockHash)
		if err != nil {
			return err
		}
	}
	return s.ChainStore.Delete(hash)
}

// NewBlock creates a new block for the provided chain and
// guarantees it as a valid next block on the chain.
func (s *Service) NewBlock(c *Chain, data interface{}) (*Block, error) {
	if s.blocks == nil {
		s.blocks = make(map[Hash]*Block)
	}
	for {
		b := c.NewBlock(data)
		err := c.AddBlock(b)
		if err == ErrInvalidPrevHash {
			continue
		} else if err != nil {
			return nil, err
		}
		s.blocks[b.Hash] = b
		return b, nil
	}
}

// NewBranch creates a branch off of a chain's block.
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

// NewChain creates new chain with a genesis block.
func (s *Service) NewChain() (*Chain, error) {
	if s.chains == nil {
		s.chains = make(map[ChainHash]*Chain)
	}
	c, err := NewChain()
	if err != nil {
		return nil, err
	}
	b := c.NewBlock(nil)
	err = c.AddBlock(b)
	if err != nil {
		return nil, err
	}
	err = s.WriteBlock(b)
	if err != nil {
		return nil, err
	}
	s.chains[c.Hash] = c
	return c, nil
}

// WriteBlock writes a block to the block store.
func (s *Service) WriteBlock(b *Block) error {
	if s.BlockStore == nil {
		return errors.New("no BlockStore provided to the service")
	}
	return s.BlockStore.Write(b)
}

// WriteChain writes a chain to the chain store
// along with any of its blocks still in memory.
func (s *Service) WriteChain(c *Chain) error {
	if s.ChainStore == nil {
		return errors.New("no ChainStore provided to the service")
	}
	err := s.ChainStore.Write(c)
	if err != nil {
		return err
	}
	if s.chains != nil {
		delete(s.chains, c.Hash)
	}
	for blockHash := range c.Blocks {
		if s.blocks[blockHash] != nil {
			err := s.WriteBlock(s.blocks[blockHash])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// NewService creates a new block service instance.
func NewService(blockStore Store, chainStore ChainStore) *Service {
	return &Service{
		BlockStore: blockStore,
		ChainStore: chainStore,

		blocks: make(map[Hash]*Block),
		chains: make(map[ChainHash]*Chain),
	}
}
