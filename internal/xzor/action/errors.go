package action

import "errors"

// ErrDuplicateAction indicates that the action service received a duplicate action.
var ErrDuplicateAction = errors.New("duplicate action received")
