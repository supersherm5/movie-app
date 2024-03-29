package grpc

import (
	"context"
	"errors"
	"github.com/supersherm5/movie-app/gen"
	"github.com/supersherm5/movie-app/metadata/internal/controller/metadata"
	"github.com/supersherm5/movie-app/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a movie metadata gRPC handler
type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

// New creates a new movie metadata gRPC handler
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMetadata return movie metadata by id.
func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "movie id is required")
	}
	m, err := h.ctrl.Get(ctx, req.MovieId)

	if err != nil && errors.Is(err, metadata.ErrMetadataNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	return &gen.GetMetadataResponse{
		Metadata: model.MetadataToProto(m),
	}, nil
}
