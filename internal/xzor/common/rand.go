package common

import "crypto/rand"

// NewRandomBytes creates a new byte slice of the provided length.
func NewRandomBytes(len int) ([]byte, error) {
	rb := make([]byte, len)
	_, err := rand.Read(rb)
	return rb, err
}

// NewRandomHash generates a random string hash at the specified length.
func NewRandomHash(len int) (string, error) {
	rb, err := NewRandomBytes(len)
	if err != nil {
		return "", err
	}
	return NewHash(rb)
}
