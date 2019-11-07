package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xzor-dev/xzor/cmd/xzor-sim/jobs"
	"github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"
	"github.com/xzor-dev/xzor/internal/module/messenger"
	msg_command "github.com/xzor-dev/xzor/internal/module/messenger/command"
	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/command"
	"github.com/xzor-dev/xzor/internal/xzor/instance"
	"github.com/xzor-dev/xzor/internal/xzor/network"
	"github.com/xzor-dev/xzor/internal/xzor/storage"
	storage_file "github.com/xzor-dev/xzor/internal/xzor/storage/file"
	storage_json "github.com/xzor-dev/xzor/internal/xzor/storage/json"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	builder := simulator.NewNetworkWebBuilder(config.Network.TotalNodes, config.Network.MaxConnectionsPerNode)
	nodes, err := builder.Build()
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	instances := make([]*instance.Instance, len(nodes))
	for i, node := range nodes {
		inst, err := newInstance(i, dir+"/testdata", node)
		if err != nil {
			panic(err)
		}
		instances[i] = inst
	}

	n, err := simulator.NewNetwork(nodes)
	if err != nil {
		panic(err)
	}

	sim := simulator.New(config.Simulator, n)

	fmt.Printf("simulator started with %d nodes\n", config.Network.TotalNodes)

	for {
		for jobID, jobConfig := range config.Jobs {
			if !jobConfig.Enabled {
				continue
			}
			go runJob(sim, jobID)
		}
		time.Sleep(time.Second * 10)
	}
}

func loadConfig() (*Config, error) {
	return &Config{
		Jobs: map[simulator.JobID]*JobConfig{
			(&jobs.MessengerJob{}).JobID(): &JobConfig{
				Enabled: true,
			},
			(&jobs.TestJob{}).JobID(): &JobConfig{
				Enabled: false,
			},
		},
		Network: &NetworkConfig{
			TotalNodes:            32,
			MaxConnectionsPerNode: 4,
		},
		Simulator: &simulator.Config{},
	}, nil
}

func newInstance(index int, dataDir string, node *network.Node) (*instance.Instance, error) {
	msgRecordDir := fmt.Sprintf("%s/instance-%d/messenger", dataDir, index)
	msgRecordStore := storage_file.NewRecordStore(msgRecordDir)
	msgStorage := storage.NewService(&storage_json.EncodeDecoder{}, msgRecordStore)
	msgService := messenger.NewService(msgStorage)
	msgCommands := msg_command.Commands(msgService)
	msgMod := messenger.NewModule(msgService, msgCommands)
	actionService := action.NewService([]command.Provider{msgMod})

	inst := instance.New(actionService, node, nil)
	err := inst.Start()
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func runJob(sim *simulator.Simulator, jobID simulator.JobID) {
	fmt.Printf("running job: %s\n", jobID)

	f, err := jobs.Registry.Get(jobID)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	j, err := f.NewJob()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	res, err := sim.RunJob(j)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	desc := `
-----------------------
JOB: %s
-----------------------
SUCCESSFUL ..... %s
EXECUTIONS ..... %d
ERRORS ......... %d
-----------------------
`

	success := "YES"
	if res.Failed {
		success = "NO"
	}

	fmt.Printf(desc, jobID, success, res.TotalExecutions, len(res.Errors))
	for i, err := range res.Errors {
		fmt.Printf("ERROR #%d: %v\n", i, err)
	}
	if len(res.Errors) > 0 {
		fmt.Print("-----------------------\n")
	}
}
