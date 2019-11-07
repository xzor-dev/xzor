package main

import (
	"github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"
)

// Config contains options for the simulator executable.
type Config struct {
	Jobs      map[simulator.JobID]*JobConfig `json:"jobs"`
	Network   *NetworkConfig                 `json:"network"`
	Simulator *simulator.Config              `json:"simulator"`
}

// NewConfig creates a new config instance.
func NewConfig() *Config {
	return &Config{}
}

// NetworkConfig defines options for the simulated network.
type NetworkConfig struct {
	TotalNodes            int `json:"totalNodes"`
	MaxConnectionsPerNode int `json:"maxConnectionsPerNode"`
}

// JobConfig contains options for individual jobs.
type JobConfig struct {
	Enabled bool        `json:"enabled"`
	Config  interface{} `json:"config"`
}
