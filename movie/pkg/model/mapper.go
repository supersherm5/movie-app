package model

import (
	"github.com/supersherm5/movie-app/gen"
	MetadataModel "github.com/supersherm5/movie-app/metadata/pkg/model"
)

// MovieDetailsToProto converts a movie details struct to a movie details proto
func MovieDetailsToProto(details *MovieDetails) *gen.MovieDetails {
	return &gen.MovieDetails{
		Metadata: &gen.Metadata{
			Id:          details.Metadata.ID,
			Title:       details.Metadata.Title,
			Description: details.Metadata.Description,
			Director:    details.Metadata.Director,
		},
		Rating: float32(details.Rating),
	}
}

// MovieDetailsFromProto converts a movie details proto to a movie details struct
func MovieDetailsFromProto(details *gen.MovieDetails) *MovieDetails {
	return &MovieDetails{
		Metadata: MetadataModel.Metadata{
			ID:          details.Metadata.Id,
			Title:       details.Metadata.Title,
			Description: details.Metadata.Description,
			Director:    details.Metadata.Director,
		},
		Rating: float64(details.Rating),
	}
}
