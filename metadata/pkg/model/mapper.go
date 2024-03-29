package model

import "github.com/supersherm5/movie-app/gen"

// MetadataToProto converts a metadata struct to a metadata proto
func MetadataToProto(metadata *Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id:          metadata.ID,
		Title:       metadata.Title,
		Description: metadata.Description,
		Director:    metadata.Director,
	}
}

// MetadataFromProto converts a metadata proto to a metadata struct
func MetadataFromProto(metadata *gen.Metadata) *Metadata {
	return &Metadata{
		ID:          metadata.Id,
		Title:       metadata.Title,
		Description: metadata.Description,
		Director:    metadata.Director,
	}
}
