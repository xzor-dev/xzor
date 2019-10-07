package simulator

import "errors"

// ErrInvalidNodeIndex is returned when an invalid node index is requested.
var ErrInvalidNodeIndex = errors.New("invalid node index")

// ErrNetworkSizeTooSmall indicates that the number of connections per node is not smaller than the network.
var ErrNetworkSizeTooSmall = errors.New("network size must be larger than connections per node")

// ErrNoConfig indicates that the simulator was not provided a config.
var ErrNoConfig = errors.New("no config found")

// ErrNoJobs is returned when the simulator is run without jobs.
var ErrNoJobs = errors.New("no jobs found")
