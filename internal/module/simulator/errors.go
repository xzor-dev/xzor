package simulator

import "errors"

// ErrNoConfig indicates that the simulator was not provided a config.
var ErrNoConfig = errors.New("no config found")

// ErrNoJobs is returned when the simulator is run without jobs.
var ErrNoJobs = errors.New("no jobs found")
