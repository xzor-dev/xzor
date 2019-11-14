package main

import (
	"fmt"
	"io"
	"log"
	"os"

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

var _ io.Writer = &logWriter{}

type logWriter struct{}

func (w *logWriter) Write(p []byte) (int, error) {
	return 0, nil
}

func main() {
	log.SetOutput(&logWriter{})

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
		inst, err := newInstance(i, dir+"/simdata", node)
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
	jobCount := 0
	for _, jc := range config.Jobs {
		if jc.Enabled {
			jobCount++
		}
	}

	fmt.Printf("# SIMULATOR STARTED\n")
	fmt.Printf("# TOTAL NODES ..... %d\n", config.Network.TotalNodes)
	fmt.Printf("# TOTAL JOBS ...... %d\n", jobCount)

	results := make(chan *simulator.JobResult, jobCount)
	for jobID, jobConfig := range config.Jobs {
		if !jobConfig.Enabled {
			continue
		}
		log.Printf("running job: %s", jobID)
		go func(jobID simulator.JobID) {
			res, err := runJob(sim, jobID)
			if err != nil {
				log.Printf("failed to run job %s: %v", jobID, err)
			}
			results <- res
		}(jobID)
	}

	for i := 0; i < jobCount; i++ {
		res := <-results
		if res == nil {
			continue
		}
		printJobResult(res)
	}

	os.RemoveAll(dir + "/simdata")
}

func loadConfig() (*Config, error) {
	return &Config{
		Jobs: map[simulator.JobID]*JobConfig{
			(&jobs.MessengerJob{}).JobID(): &JobConfig{
				Enabled: true,
			},
			(&jobs.TestJob{}).JobID(): &JobConfig{
				Enabled: true,
			},
		},
		Network: &NetworkConfig{
			TotalNodes:            64,
			MaxConnectionsPerNode: 6,
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

func printJobResult(res *simulator.JobResult) {
	success := "YES"
	if res.Failed {
		success = "NO"
	}

	dur := res.EndTime.Sub(res.StartTime)

	fmt.Printf("\n")
	fmt.Printf("-----------------\n")
	fmt.Printf("| JOB ........... %s\n", res.JobID)
	fmt.Printf("| DURATION ...... %s\n", dur)
	fmt.Printf("| SUCCESS ....... %s\n", success)

	fmt.Printf("| ERRORS ........ %d\n", len(res.Errors))
	for i, err := range res.Errors {
		fmt.Printf("|     #%d: %s\n", i+1, err)
	}

	fmt.Printf("| ACTIONS ....... %d\n", len(res.Actions))
	for i, a := range res.Actions {
		hash := []rune(a.Hash)
		hashSub := string(hash[0:5]) + "..." + string(hash[len(hash)-5:])
		fmt.Printf("|     #%d: %s\n", i+1, hashSub)
	}

	fmt.Printf("-----------------\n")
}

func runJob(sim *simulator.Simulator, jobID simulator.JobID) (*simulator.JobResult, error) {
	f, err := jobs.Registry.Get(jobID)
	if err != nil {
		return nil, err
	}
	j, err := f.NewJob()
	if err != nil {
		return nil, err
	}
	return sim.RunJob(j)
}
