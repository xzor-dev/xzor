package file

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/xzor-dev/xzor/internal/xzor/block"
)

var _ block.Store = &BlockStore{}

// BlockStore provides reading and writing of blocks to the file system.
type BlockStore struct {
	RootDir string
}

func (s *BlockStore) filename(hash block.Hash) string {
	return s.RootDir + "/" + string(hash)
}

// Delete removes the block's data file.
func (s *BlockStore) Delete(hash block.Hash) error {
	return os.Remove(s.filename(hash))
}

// Read gets a block using its hash.
func (s *BlockStore) Read(hash block.Hash) (*block.Block, error) {
	filename := s.filename(hash)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	b := &block.Block{}
	err = json.Unmarshal(data, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Write a block to the file system.
func (s *BlockStore) Write(b *block.Block) error {
	err := os.MkdirAll(s.RootDir, 0666)
	if err != nil {
		return err
	}

	filename := s.filename(b.Hash)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(b)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// NewBlockStore creates a new file-based block store instance.
func NewBlockStore(rootDir string) *BlockStore {
	return &BlockStore{
		RootDir: rootDir,
	}
}
