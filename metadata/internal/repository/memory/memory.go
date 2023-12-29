package memory

import (
	"context"
	"sync"

	"github.com/supersherm5/movie-app/metadata/internal/repository"
	"github.com/supersherm5/movie-app/metadata/pkg/model"
)

// Repository defines a memory movie metadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New Repository returns a new instance of a memory movie metadata repository.
func New() *Repository {
	return &Repository{
		data: make(map[string]*model.Metadata),
	}
}

// Get returns a movie metadata from the repository by id
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()

	if m, ok := r.data[id]; !ok {
		return m, nil
	}

	return nil, repository.ErrMetadataNotFound
}

// Put adds movie metadata for a given movie id to the repository
func (r *Repository) Put(_ context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
