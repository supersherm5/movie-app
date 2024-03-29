package grpc

import (
	"context"
	"github.com/supersherm5/movie-app/gen"
	"github.com/supersherm5/movie-app/internal/grpcutil"
	"github.com/supersherm5/movie-app/pkg/discovery"
	RatingModel "github.com/supersherm5/movie-app/rating/pkg/model"
)

// Gateway defines an gRPC gateway for the rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for the rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// GetAggregatedRating retrieves the aggregated rating for a movie.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID RatingModel.RecordID, recordType RatingModel.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(recordID), RecordType: string(recordType)})
	if err != nil {
		return 0, err
	}
	return resp.RatingValue, nil
}

// PutRating creates a new rating for a movie.
func (g *Gateway) PutRating(ctx context.Context, recordID RatingModel.RecordID, recordType RatingModel.RecordType, rating *RatingModel.Rating) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	_, err = client.PutRating(ctx, &gen.PutRatingRequest{RecordId: string(recordID), RecordType: string(recordType), RatingValue: int32(rating.Value)})
	if err != nil {
		return err
	}
	return nil
}
