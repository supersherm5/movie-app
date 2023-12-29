package repository

import "errors"

// ErrMetadataNotFound is returned when the metadata is not found.
var ErrMetadataNotFound = errors.New("metadata not found")
