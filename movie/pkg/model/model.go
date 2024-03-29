package model

import "github.com/supersherm5/movie-app/metadata/pkg/model"

// MovieDetails includes movie metadata and its rating.
type MovieDetails struct {
	Rating   float64        `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}
