package resource

import "errors"

// ErrInvalidProvider indicates that a provider doesn't exist.
var ErrInvalidProvider = errors.New("invalid provider")

// ErrInvalidResourceName indicates that a resource name doesn't exist for a provider.
var ErrInvalidResourceName = errors.New("invalid resource name")
