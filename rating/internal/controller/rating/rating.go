package rating

import (
	"context"
	"errors"
	"github.com/supersherm5/movie-app/rating/internal/repository"
	"github.com/supersherm5/movie-app/rating/pkg/model"
)

// ErrRatingTypeNotFound is returned when the rating record id is not found.
var ErrRatingTypeNotFound = errors.New("rating type id not found")

// ErrRatingRecordIDNotFound is returned when the rating record id is not found.
var ErrRatingRecordIDNotFound = errors.New("rating record id not found")

// ErrNoRatingExists is returned when there are ratings, but no value has been created
var ErrNoRatingExists = errors.New("no rating exists for this record")

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}
type ratingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}
// Controller defines a rating service controller.
type Controller struct {
	repo ratingRepository
	ingester	ratingIngester

}

// New creates a rating service controller.
// New creates a rating service
func New(repo ratingRepository, ingester ratingIngester) *Controller {
	return &Controller{repo: repo, ingester: ingester}
}

// GetAggregatedRating returns the aggregated rating for an ErrNotRatingExists if no rating exists for the record.
func (c *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)

	if err != nil && errors.Is(err, repository.ErrRatingTypeNotFound) {
		return 0, ErrRatingTypeNotFound
	} else if err != nil && errors.Is(err, repository.ErrRatingRecordIDNotFound) {
		return 0, ErrRatingRecordIDNotFound
	} else if err != nil && len(ratings) == 0 {
		return 0, ErrNoRatingExists
	}

	// Calculate the average rating
	sum := float64(0)
	for _, rating := range ratings {
		sum += float64(rating.Value)
	}

	return sum / float64(len(ratings)), nil
}

func (c *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, recordID, recordType, rating)
}


// StartIngestion starts the ingestion of rating events.
func (c *Controller) StartIngestion(ctx context.Context) error {
	ch, err := c.ingester.Ingest(ctx)
	if err != nil {
		return err
	}

	for e := range ch {
		if err := c.repo.Put(ctx, e.RecordID, e.RecordType, &model.Rating{UserID: e.UserID, Value: e.Value }); err != nil {
			return err
		}
	}

	return nil
}