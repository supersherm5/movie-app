package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/supersherm5/movie-app/gen"
	"github.com/supersherm5/movie-app/rating/internal/controller/rating"
	"github.com/supersherm5/movie-app/rating/internal/repository/memory"
	"github.com/supersherm5/movie-app/rating/pkg/model"
)

// Handler defines a movie rating gRPC handler
type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *rating.Controller
}

// New creates a new movie rating gRPC handler
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetAggregatedRating returns the aggregated rating for a
// record.
func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	v, err := h.ctrl.GetAggregatedRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, rating.ErrRatingTypeNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}

// PutRating adds a rating to a record.
func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" || req.RatingValue < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req, empty id, or negative rating value")
	}
	err := h.ctrl.PutRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), memory.RatingModelFromProto(req))

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.PutRatingResponse{}, nil
}
