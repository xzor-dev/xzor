package action

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/xzor/block"
)

type Aggregator struct {
	BlockService *block.Service
	Chain        *block.Chain

	actions []*Action
}

func (a *Aggregator) Clear() {
	a.actions = make([]*Action, 0)
}

func (a *Aggregator) GenerateBlock() (*block.Block, error) {
	if a.BlockService == nil {
		return nil, errors.New("no block service provided to the aggregator")
	}
	if a.Chain == nil {
		return nil, errors.New("no chain provided to the aggregator")
	}
	b, err := a.BlockService.NewBlock(a.Chain, a.actions)
	if err != nil {
		return nil, err
	}
	a.Clear()
	return b, nil
}

func (a *Aggregator) Push(action *Action) {
	if a.actions == nil {
		a.actions = make([]*Action, 0)
	}
	a.actions = append(a.actions, action)
}

func NewAggregator(blockService *block.Service, chain *block.Chain) *Aggregator {
	return &Aggregator{
		BlockService: blockService,
		Chain:        chain,

		actions: make([]*Action, 0),
	}
}
