package action_test

import (
	"os"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/block"
	"github.com/xzor-dev/xzor/internal/xzor/block/file"
)

func TestAggregator(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}

	blockStore := file.NewBlockStore(dir + "/testdata/blocks")
	chainStore := file.NewChainStore(dir + "/testdata/chains")
	blockService := block.NewService(blockStore, chainStore)
	chain, err := block.NewChain()
	if err != nil {
		t.Fatalf("%v", err)
	}

	ag := action.NewAggregator(blockService, chain)
	ag.Push(&action.Action{
		Command:         "test-command",
		CommandProvider: "test-module",
		Parameters: map[string]interface{}{
			"test": "foo",
		},
	})
	block, err := ag.GenerateBlock()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if block.Data == nil {
		t.Fatal("expected block to have action data")
	}
	if !chain.HasBlock(block.Hash) {
		t.Fatal("expected chain to have block")
	}
}
