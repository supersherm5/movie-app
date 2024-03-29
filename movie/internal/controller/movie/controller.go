package movie

import (
	"context"
	"errors"
	MetadataModel "github.com/supersherm5/movie-app/metadata/pkg/model"
	"github.com/supersherm5/movie-app/movie/internal/gateway"
	"github.com/supersherm5/movie-app/movie/pkg/model"
	RatingModel "github.com/supersherm5/movie-app/rating/pkg/model"
)

// ErrNotFound is returned when movie metadata is not found.
var ErrNotFound = errors.New("movie metadata not found")

type RatingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID RatingModel.RecordID, recordType RatingModel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID RatingModel.RecordID, recordType RatingModel.RecordType, rating *RatingModel.Rating) error
}

type MetadataGateway interface {
	Get(ctx context.Context, id string) (*MetadataModel.Metadata, error)
}

// Controller defines a movie controller.
type Controller struct {
	ratingGateway   RatingGateway
	metadataGateway MetadataGateway
}

// New creates a new movie controller.
func New(ratingGateway RatingGateway, metadataGateway MetadataGateway) *Controller {
	return &Controller{ratingGateway, metadataGateway}
}

// Get returns the movie details including the aggregated rating and movie metadata.
func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gateway.ErrServiceNotReachable) {
			return nil, gateway.ErrServiceNotReachable
		} else if errors.Is(err, gateway.ErrNotFound) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, RatingModel.RecordID(id), RatingModel.RecordTypeMovie)
	if err != nil {
		if errors.Is(err, gateway.ErrServiceNotReachable) {
			return nil, gateway.ErrServiceNotReachable
		} else if errors.Is(err, gateway.ErrNotFound) {
			return details, nil
		} else {
			return nil, err
		}
	}

	details.Rating = rating
	return details, nil
}

func (c *Controller) PutRating(ctx context.Context, id, recordType string, rating *RatingModel.Rating) error {
	return c.ratingGateway.PutRating(ctx, RatingModel.RecordID(id), RatingModel.RecordType(recordType), rating)
}
