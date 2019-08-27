package file

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/xzor-dev/xzor/internal/xzor/storage"
)

var _ storage.ChainStore = &ChainStore{}

// ChainStore provides reading and writing of chains to the file system.
type ChainStore struct {
	RootDir string
}

func (s *ChainStore) filename(hash storage.ChainHash) string {
	return s.RootDir + "/" + string(hash)
}

// Delete removes the chain's data file.
func (s *ChainStore) Delete(hash storage.ChainHash) error {
	return os.Remove(s.filename(hash))
}

// Read chain's data using its hash.
func (s *ChainStore) Read(hash storage.ChainHash) (*storage.Chain, error) {
	filename := s.filename(hash)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &storage.Chain{}
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Write a chain to the file system.
func (s *ChainStore) Write(c *storage.Chain) error {
	err := os.MkdirAll(s.RootDir, 0666)
	if err != nil {
		return err
	}

	filename := s.filename(c.Hash)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
