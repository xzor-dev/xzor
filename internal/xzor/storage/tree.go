package storage

// Branch holds a chain identifed by a BranchHash.
type Branch struct {
	Chain *Chain
	Hash  BranchHash
}

// BranchHash is a unique string assigned to branches.
type BranchHash string

// NewBranchHash generates a new unique branch hash.
func NewBranchHash() (BranchHash, error) {
	return "", nil
}

// Tree holds multiple branches identified by a TreeHash.
type Tree struct {
	Branches map[BranchHash]*Branch
	Hash     TreeHash
}

// TreeHash is a unique string assigned to trees.
type TreeHash string

// NewTreeHash generates a new unique tree hash.
func NewTreeHash() (TreeHash, error) {
	return "", nil
}
