package memory

import (
	"github.com/supersherm5/movie-app/gen"
	RatingModel "github.com/supersherm5/movie-app/rating/pkg/model"
)

func RatingResponseToProto(rating *RatingModel.Rating) *gen.GetAggregatedRatingResponse {
	return &gen.GetAggregatedRatingResponse{
		RatingValue: float64(rating.Value),
	}
}

func RatingModelFromProto(rating *gen.PutRatingRequest) *RatingModel.Rating {
	return &RatingModel.Rating{
		RecordID:   RatingModel.RecordID(rating.RecordId),
		RecordType: RatingModel.RecordType(rating.RecordType),
		UserID:     RatingModel.UserID(rating.UserId),
		Value:      RatingModel.RatingValue(rating.RatingValue),
	}
}
