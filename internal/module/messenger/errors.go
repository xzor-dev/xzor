package messenger

import "errors"

// ErrInvalidBoardHash occurs when a board hash is missing or invalid.
var ErrInvalidBoardHash = errors.New("invalid board hash")

// ErrNoStorage indicates that a storage service was not found.
var ErrNoStorage = errors.New("no storage service found")
