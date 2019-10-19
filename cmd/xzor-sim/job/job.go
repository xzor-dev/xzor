package job

import "github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"

type ID string

type Factory interface {
	ID() ID
	NewJob() (simulator.Job, error)
}

var Jobs = make(map[ID]Factory)
