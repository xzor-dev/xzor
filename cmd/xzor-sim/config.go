package main

import (
	"github.com/xzor-dev/xzor/cmd/xzor-sim/job"
	"github.com/xzor-dev/xzor/cmd/xzor-sim/simulator"
)

// Config contains options for the simulator executable.
type Config struct {
	Jobs     map[job.ID]bool
	SimCofig *simulator.Config
}

// NewConfig creates a new config instance.
func NewConfig() *Config {
	return &Config{}
}
