package metadata

import (
	"context"
	"errors"

	"github.com/supersherm5/movie-app/metadata/internal/repository"
	"github.com/supersherm5/movie-app/metadata/pkg/model"
)

// ErrMetadataNotFound is returned when the metadata is not found.
var ErrMetadataNotFound = errors.New("metadata not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

// Controller defines a movie metadata service controller.
type Controller struct {
	repo metadataRepository
}

// New creates a metadata service controller.
func New(repo metadataRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}

// Get returns a movie metadata by id.
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrMetadataNotFound) {
		return nil, ErrMetadataNotFound
	}
	return res, nil
}
