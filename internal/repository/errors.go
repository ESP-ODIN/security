package repository

import "errors"

// ErrNotFound is returned when a version or a pipeline run does not exist.
var ErrNotFound = errors.New("not found")
