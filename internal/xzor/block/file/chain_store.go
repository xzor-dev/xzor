package file

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/xzor-dev/xzor/internal/xzor/block"
)

var _ block.ChainStore = &ChainStore{}

// ChainStore provides reading and writing of chains to the file system.
type ChainStore struct {
	RootDir string
}

func (s *ChainStore) filename(hash block.ChainHash) string {
	return s.RootDir + "/" + string(hash)
}

// Delete removes the chain's data file.
func (s *ChainStore) Delete(hash block.ChainHash) error {
	return os.Remove(s.filename(hash))
}

// Read chain's data using its hash.
func (s *ChainStore) Read(hash block.ChainHash) (*block.Chain, error) {
	filename := s.filename(hash)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &block.Chain{}
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Write a chain to the file system.
func (s *ChainStore) Write(c *block.Chain) error {
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

// NewChainStore creates a new file-based chain store instance.
func NewChainStore(rootDir string) *ChainStore {
	return &ChainStore{
		RootDir: rootDir,
	}
}
