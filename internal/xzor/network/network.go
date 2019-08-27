package network

import (
	"crypto/sha256"
	"net"
)

type Connection interface {
	Connect() (net.Conn, error)
}

type DataHandler interface {
	HandleData([]byte) error
}

type Listener interface {
	Listen() (net.Listener, error)
}

type Request struct {
	Data []byte
	Hash []byte
}

func (r *Request) GenerateHash() error {
	h := sha256.Sum256(r.Data)
	r.Hash = h[:32]
	return nil
}
