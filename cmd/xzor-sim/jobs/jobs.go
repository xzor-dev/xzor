package jobs

import (
	"errors"
	"log"

	"github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"
)

type Factory interface {
	JobID() simulator.JobID
	NewJob() (simulator.Job, error)
}

type FactoryMap map[simulator.JobID]Factory

func (fm FactoryMap) Get(id simulator.JobID) (Factory, error) {
	if fm[id] == nil {
		return nil, errors.New("invalid job ID")
	}
	return fm[id], nil
}

func (fm FactoryMap) Set(f Factory) {
	if fm[f.JobID()] != nil {
		log.Printf("replacing existing factory: %s", f.JobID())
		delete(fm, f.JobID())
	}
	fm[f.JobID()] = f
}

var Registry = make(FactoryMap)
