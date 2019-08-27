package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

// Block holds sequential data.
type Block struct {
	Data         []byte
	Hash         BlockHash
	Index        int
	PreviousHash BlockHash
	Timestamp    int64
}

// BlockHash is a unique string generated from a block.
type BlockHash string

// NewBlockHash generates a unique hash for a block.
func NewBlockHash(b *Block) (BlockHash, error) {
	var hash BlockHash

	record := string(b.Index) + string(b.Timestamp) + string(b.Data) + string(b.PreviousHash)
	hasher := sha256.New()
	_, err := hasher.Write([]byte(record))
	if err != nil {
		return hash, err
	}
	hashBytes := hasher.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)
	hash = BlockHash(hashHex)

	return hash, nil
}
