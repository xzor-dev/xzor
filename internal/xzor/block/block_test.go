package block_test

import (
	"log"
	"os"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/block"
	"github.com/xzor-dev/xzor/internal/xzor/block/file"
	"github.com/xzor-dev/xzor/internal/xzor/block/memory"
)

func TestChain(t *testing.T) {
	c1 := &block.Chain{}
	b1 := c1.NewBlock(nil)
	err := c1.AddBlock(b1)
	if err != nil {
		t.Fatalf("%v", err)
	}

	b2 := c1.NewBlock(nil)
	err = c1.AddBlock(b2)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if b2.Hash == b1.Hash {
		t.Fatal("block has the same hash as the previous block")
	}
	if b2.Index != 1 {
		t.Fatalf("expected block to have an index of 1, got %d", b2.Index)
	}

	b3 := c1.NewBlock(nil)
	b3.Hash = "bad_hash"
	err = c1.AddBlock(b3)
	if err == nil {
		t.Fatal("expected an error when attempting to add a bad block to the chain")
	}
}

func TestService(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}
	srv := &block.Service{
		ChainStore: &file.ChainStore{
			RootDir: dir + "/testdata",
		},
	}
	c1, err := srv.NewChain()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if c1.Hash == "" {
		t.Fatal("expected newly created chain to have a hash")
	}
	if len(c1.Blocks) != 1 {
		t.Fatalf("expected chain to have 1 block, got %d", len(c1.Blocks))
	}
	err = srv.WriteChain(c1)
	if err != nil {
		t.Fatalf("%v", err)
	}
	c1A, err := srv.ReadChain(c1.Hash)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if c1A.Hash != c1.Hash {
		t.Fatalf("expected loaded chain to have same hash: wanted %s, got %s", c1.Hash, c1A.Hash)
	}
	err = srv.DeleteChain(c1.Hash)
	if err != nil {
		t.Fatalf("%v", err)
	}
	_, err = srv.ReadChain(c1.Hash)
	if err == nil {
		t.Fatalf("expected an error when reading deleted chain")
	}
}

func TestConcurrentWrites(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}
	s := &block.Service{
		ChainStore: &file.ChainStore{
			RootDir: dir + "/testdata",
		},
	}
	c, err := s.NewChain()
	if err != nil {
		t.Fatalf("%v", err)
	}

	write := func(s *block.Service, c *block.Chain, count int) error {
		for i := 0; i < count; i++ {
			_, err := s.NewBlock(c, nil)
			if err != nil {
				return err
			}
		}
		return nil
	}

	threads := 5
	writesPerThread := 10
	expectedBlocks := threads*writesPerThread + 1
	errs := make(chan error)
	for i := 0; i < threads; i++ {
		go func(i int) {
			log.Printf("writing on thread %d", i)
			err := write(s, c, writesPerThread)
			errs <- err
		}(i)
	}
	for i := 0; i < threads; i++ {
		err := <-errs
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
	if len(c.Blocks) != expectedBlocks {
		t.Fatalf("unexpected number of blocks: wanted %d, got %d", expectedBlocks, len(c.Blocks))
	}

	err = s.WriteChain(c)
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = s.DeleteChain(c.Hash)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestBranchingChains(t *testing.T) {
	s := &block.Service{
		ChainStore: &memory.ChainStore{},
	}
	c1, err := s.NewChain()
	if err != nil {
		t.Fatalf("%v", err)
	}

	b1, err := s.NewBlock(c1, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	branch, err := s.NewBranch(c1, b1)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if c1.Branches[branch.Hash] == nil {
		t.Fatalf("expected chain to have branch")
	}
}
