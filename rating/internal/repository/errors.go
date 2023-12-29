package repository

import "errors"

// ErrRatingTypeNotFound is returned when no ratings of a specified type are found.
var ErrRatingTypeNotFound = errors.New("rating type not found")

// ErrRatingRecordIDNotFound is return when no ratings are found.
var ErrRatingRecordIDNotFound = errors.New("record id not found")

// ErrNoRatingExists is returned when there are ratings, but no value has been created
var ErrNoRatingExists = errors.New("no rating exists for this record")
