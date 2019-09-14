package common

import (
	"crypto/sha256"
	"encoding/hex"
)

// NewHash generates a sha256 sum of the provided data and converts it to a hex string.
func NewHash(data []byte) (string, error) {
	var hash string
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return hash, err
	}

	hb := hasher.Sum(nil)
	hash = hex.EncodeToString(hb)

	return hash, nil
}
