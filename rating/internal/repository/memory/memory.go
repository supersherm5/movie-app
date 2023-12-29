package memory

import (
	"context"

	"github.com/supersherm5/movie-app/rating/internal/repository"
	"github.com/supersherm5/movie-app/rating/pkg/model"
)

// Repository defines a rating repository.
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{
		make(map[model.RecordType]map[model.RecordID][]model.Rating),
	}
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrRatingTypeNotFound
	}

	ratings, ok := r.data[recordType][recordID]
	if !ok {
		return nil, repository.ErrRatingRecordIDNotFound
	}

	if len(ratings) == 0 {
		return nil, repository.ErrNoRatingExists
	}

	return ratings, nil
}

// Put adds a rating to a record in the repository.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = make(map[model.RecordID][]model.Rating)
	}

	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)

	return nil
}
