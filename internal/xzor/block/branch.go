package block

import (
	"github.com/xzor-dev/xzor/internal/xzor/common"
)

// Branch holds information on branched chains.
type Branch struct {
	FromBlock Hash
	Hash      BranchHash
	ToChain   ChainHash
}

// BranchHash is a unique string assigned to branches.
type BranchHash string

// NewBranchHash generates new branch hashes.
func NewBranchHash() (BranchHash, error) {
	var bh BranchHash

	rb, err := common.NewRandomBytes(32)
	if err != nil {
		return bh, err
	}
	hash, err := common.NewHash(rb)
	if err != nil {
		return bh, err
	}
	bh = BranchHash(hash)
	return bh, nil
}
