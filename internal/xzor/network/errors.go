package network

import "errors"

// ErrNoMessages indicates that a node has no messages to read from its queue.
var ErrNoMessages = errors.New("no messages in queue")
