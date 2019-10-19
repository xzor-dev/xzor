package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xzor-dev/xzor/cmd/xzor-sim/job"
	"github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"
)

var _ simulator.Job = &actionJob{}

type actionJob struct{}

func (j *actionJob) Config() *simulator.JobConfig {
	return &simulator.JobConfig{}
}

func (j *actionJob) Execute(p *simulator.JobParams) error {
	return nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	n, err := simulator.NewNetwork(config.SimCofig, &simulator.NetworkWebBuilder{})
	sim := simulator.New(config.SimCofig, n)
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("simulator started with %d nodes\n", config.SimCofig.TotalNodes)

	for {
		fmt.Print("\n> ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error: %v", err)
			continue
		}
		cmd = strings.Replace(cmd, "\n", "", -1)
		cmdParts := strings.Split(cmd, " ")
		if len(cmdParts) < 2 {
			fmt.Println("expecting at least 2 arguments")
			fmt.Println("usage: module action [param1, [param2, ...]]")
			continue
		}

		moduleName := cmdParts[0]
		actionName := cmdParts[1]
		params := cmdParts[2:]

		fmt.Printf("got command: %s.%s %s\n", moduleName, actionName, strings.Join(params, ", "))

		job := &actionJob{}
		res, err := sim.RunJob(job)
		if err != nil {
			fmt.Printf("job failed: %v\n", err)
			continue
		}

		fmt.Println("job result:")
		fmt.Printf("\texecutions: %d\n", res.TotalExecutions)
		fmt.Printf("\terrors: %d\n", len(res.Errors))
		fmt.Printf("\tfailed: %v\n", res.Failed)
	}
}

func loadConfig() (*Config, error) {
	return &Config{
		Jobs: map[job.ID]bool{
			job.TestJobID: true,
		},
		SimCofig: &simulator.Config{
			TotalNodes:         32,
			ConnectionsPerNode: 4,
		},
	}, nil
}

func generateJobs(c *Config) ([]simulator.Job, error) {
	jobs := make([]simulator.Job, 0)
	for id, enabled := range c.Jobs {
		if !enabled {
			continue
		}
		if job.Jobs[id] == nil {
			return nil, fmt.Errorf("invalid job ID: %s", id)
		}
		j, err := job.Jobs[id].NewJob()
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}
