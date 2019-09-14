package common

import "crypto/rand"

// NewRandomBytes creates a new byte slice of the provided length.
func NewRandomBytes(len int) ([]byte, error) {
	rb := make([]byte, len)
	_, err := rand.Read(rb)
	return rb, err
}
