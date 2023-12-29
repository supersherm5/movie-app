package gateway

import "errors"

// ErrNotFound is returned when data is not found.
var ErrNotFound = errors.New("data not found")

var ErrServiceNotReachable = errors.New("service not reachable")
