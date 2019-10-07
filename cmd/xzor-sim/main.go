package main

import (
	"fmt"

	"github.com/xzor-dev/xzor/cmd/xzor-sim/job"
	"github.com/xzor-dev/xzor/internal/module/simulator"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	n, err := simulator.NewNetwork(config.SimCofig, &simulator.NetworkWebBuilder{})
	sim := simulator.New(config.SimCofig, n)
	jobs, err := generateJobs(config)
	if err != nil {
		panic(err)
	}
	res, err := sim.Run(jobs)
	if err != nil {
		panic(err)
	}

	fmt.Println("simulation finished")
	fmt.Printf("Total Jobs Executed .... %d\n", res.TotalExecuted)
	fmt.Printf("Total Jobs Completed ... %d\n", res.TotalCompleted)
	fmt.Printf("Total Job Failures ..... %d\n", res.TotalFailed)
	for i, jRes := range res.Jobs {
		fmt.Printf("\tJob #%d -- Executions: %d, Errors: %d, Failed: %v\n", i, jRes.TotalExecutions, len(jRes.Errors), jRes.Failed)
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
