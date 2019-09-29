package job

import "github.com/xzor-dev/xzor/internal/module/simulator"

type ID string

type Factory interface {
	ID() ID
	NewJob() (simulator.Job, error)
}

var Jobs = make(map[ID]Factory)
