package storage_test

import (
	"log"
	"os"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/storage"
	"github.com/xzor-dev/xzor/internal/xzor/storage/file"
)

func TestChain(t *testing.T) {
	c1 := &storage.Chain{}
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
	srv := &storage.Service{
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
	s := &storage.Service{}
	c, err := s.NewChain()
	if err != nil {
		t.Fatalf("%v", err)
	}

	write := func(t *testing.T, c *storage.Chain, count int) {
		for i := 0; i < count; i++ {
			b := c.NewBlock(nil)
			err := c.AddBlock(b)
			if err != nil {
				t.Fatalf("%v", err)
			}
		}
	}

	threads := 5
	writesPerThread := 100
	expectedBlocks := threads*writesPerThread + 1
	done := make(chan bool)
	for i := 0; i < threads; i++ {
		go func(i int) {
			log.Printf("writing on thread %d", i)
			write(t, c, writesPerThread)
			done <- true
		}(i)
	}
	for i := 0; i < threads; i++ {
		<-done
	}
	if len(c.Blocks) != expectedBlocks {
		t.Fatalf("unexpected number of blocks: wanted %d, got %d", expectedBlocks, len(c.Blocks))
	}
}
