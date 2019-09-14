package block

import (
	"sync"
	"time"

	"github.com/xzor-dev/xzor/internal/xzor/common"
)

// Chain holds a set of blocks and enforces their ordering.
type Chain struct {
	Blocks   map[Hash]Index
	Branches map[BranchHash]*Branch
	Hash     ChainHash
	LastHash Hash

	mux sync.Mutex
}

// AddBlock adds a new block to the chain.
func (c *Chain) AddBlock(b *Block) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.Blocks == nil {
		c.Blocks = make(map[Hash]Index)
	}

	hash, err := NewHash(b)
	if err != nil {
		return err
	}

	if b.Hash == "" {
		b.Hash = hash
	} else if hash != b.Hash {
		return ErrInvalidHash
	}

	if c.LastHash != "" {
		lastIndex := c.Blocks[c.LastHash]

		if b.PreviousHash != c.LastHash {
			return ErrInvalidPrevHash
		}
		if b.Index != lastIndex+1 {
			return ErrInvalidIndex
		}
	}

	c.Blocks[b.Hash] = b.Index
	c.LastHash = b.Hash

	return nil
}

// NewBlock creates a new block with all required properties pre-populated.
func (c *Chain) NewBlock(data []byte) *Block {
	c.mux.Lock()
	defer c.mux.Unlock()

	index := 0
	if c.LastHash != "" {
		lastIndex := c.Blocks[c.LastHash]
		index = int(lastIndex) + 1
	}

	return &Block{
		Data:         data,
		Index:        Index(index),
		PreviousHash: c.LastHash,
		Timestamp:    time.Now().Unix(),
	}
}

// NewBranch creates a new branch off of the provided block to the provided chain.
func (c *Chain) NewBranch(fromBlock *Block, toChain *Chain) (*Branch, error) {
	if c.Branches == nil {
		c.Branches = make(map[BranchHash]*Branch)
	}
	if c.Blocks == nil {
		return nil, ErrEmptyChain
	}
	if c.Blocks[fromBlock.Hash] == 0 {
		return nil, ErrInvalidHash
	}

	hash, err := NewBranchHash()
	if err != nil {
		return nil, err
	}
	branch := &Branch{
		FromBlock: fromBlock.Hash,
		Hash:      hash,
		ToChain:   toChain.Hash,
	}
	c.Branches[hash] = branch
	return branch, nil
}

// ChainHash is a unique string assigned to chains.
type ChainHash string

// NewChainHash generates a new unique chain hash.
func NewChainHash() (ChainHash, error) {
	var ch ChainHash

	rb, err := common.NewRandomBytes(32)
	if err != nil {
		return ch, err
	}

	hash, err := common.NewHash(rb)
	if err != nil {
		return ch, err
	}
	ch = ChainHash(hash)

	return ch, nil
}

// ChainStore handles storage operations for chains.
type ChainStore interface {
	Delete(ChainHash) error
	Read(ChainHash) (*Chain, error)
	Write(*Chain) error
}
