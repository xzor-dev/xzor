package jobs

import (
	"log"

	"github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"
	"github.com/xzor-dev/xzor/internal/xzor/action"
)

func init() {
	Registry.Set(&MessengerJobFactory{})
}

const messengerJobID = "messenger"

var _ simulator.Job = &MessengerJob{}

type MessengerJob struct{}

func (m *MessengerJob) Config() *simulator.JobConfig {
	return &simulator.JobConfig{}
}

func (m *MessengerJob) Execute(p *simulator.JobParams) (*action.Action, error) {
	a, err := action.New("messenger", "create-board", map[string]interface{}{
		"title": "Test Board",
	})
	log.Printf("new action created: %s at %d", a.Hash, a.Timestamp)
	return a, err
}

func (m *MessengerJob) JobID() simulator.JobID {
	return messengerJobID
}

var _ Factory = &MessengerJobFactory{}

type MessengerJobFactory struct{}

func (f *MessengerJobFactory) JobID() simulator.JobID {
	return messengerJobID
}

func (f *MessengerJobFactory) NewJob() (simulator.Job, error) {
	return &MessengerJob{}, nil
}
