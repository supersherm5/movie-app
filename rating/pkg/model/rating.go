package model

// RecordID defines a record id. Together with RecordType it uniquely identifies records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID identifies records across all types.
type RecordType string

// UserID defines a user id.
type UserID string

// RatingValue RatingValue Type defines a rating value.
type RatingValue float64

// RecordTypeMovie Existing record types.
const (
	RecordTypeMovie RecordType = "movie"
)

// Rating defines an individual rating created by a user for some record.
type Rating struct {
	RecordID   RecordID    `json:"record_id"`
	RecordType RecordType  `json:"record_type"`
	UserID     UserID      `json:"user_id"`
	Value      RatingValue `json:"value"`
}
