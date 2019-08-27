package storage

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// Chain holds a set of blocks and enforces their ordering.
type Chain struct {
	Blocks []*Block
	Hash   ChainHash

	mux sync.Mutex
}

// AddBlock adds a new block to the chain.
func (c *Chain) AddBlock(b *Block) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.Blocks == nil {
		c.Blocks = make([]*Block, 0)
	}

	hash, err := NewBlockHash(b)
	if err != nil {
		return err
	}

	if b.Hash == "" {
		b.Hash = hash
	} else if hash != b.Hash {
		return errors.New("invalid block hash")
	}

	c.Blocks = append(c.Blocks, b)
	return nil
}

// NewBlock creates a new block with all required properties pre-populated.
func (c *Chain) NewBlock(data []byte) *Block {
	c.mux.Lock()
	defer c.mux.Unlock()

	var prevHash BlockHash

	index := 0
	if c.Blocks != nil && len(c.Blocks) > 0 {
		prevBlock := c.Blocks[len(c.Blocks)-1]
		prevHash = prevBlock.Hash
		index = prevBlock.Index + 1
	}

	return &Block{
		Data:         data,
		Index:        index,
		PreviousHash: prevHash,
		Timestamp:    time.Now().Unix(),
	}
}

// ChainHash is a unique string assigned to chains.
type ChainHash string

// NewChainHash generates a new unique chain hash.
func NewChainHash() (ChainHash, error) {
	var hash ChainHash

	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil {
		return hash, err
	}

	hasher := sha256.New()
	_, err = hasher.Write(rb)
	if err != nil {
		return hash, err
	}

	hb := hasher.Sum(nil)
	hh := hex.EncodeToString(hb)
	hash = ChainHash(hh)

	return hash, nil
}

// ChainStore handles storage operations for chains.
type ChainStore interface {
	Delete(ChainHash) error
	Read(ChainHash) (*Chain, error)
	Write(*Chain) error
}
