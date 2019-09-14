package block

import "errors"

// ErrEmptyChain indicates when a chain is empty.
var ErrEmptyChain = errors.New("empty chain")

// ErrInvalidHash occurs when a block's hash is found to be invalid.
var ErrInvalidHash = errors.New("invalid block hash")

// ErrInvalidIndex occurs when a block's index is not correct for its chain.
var ErrInvalidIndex = errors.New("invalid block index")

// ErrInvalidPrevHash occurs when a block's previous hash does not match the chain's last block hash.
var ErrInvalidPrevHash = errors.New("invalid previous block hash")
