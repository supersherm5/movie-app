package http

import (
	"encoding/json"
	"errors"
	"github.com/supersherm5/movie-app/rating/internal/controller/rating"
	"github.com/supersherm5/movie-app/rating/pkg/model"
	"log"
	"net/http"
	"strconv"
)

// Handler defines a rating service http handler.
type Handler struct {
	ctrl *rating.Controller
}

// New creates a new rating service HTTP handler.
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

// Handle handles PUT and GET /rating requests.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	recordID := model.RecordID(r.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(r.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAggregatedRating(r.Context(), recordID, recordType)
		if err != nil {
			if errors.Is(err, rating.ErrRatingTypeNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if errors.Is(err, rating.ErrRatingRecordIDNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if errors.Is(err, rating.ErrNoRatingExists) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("Response encoding error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		userID := model.UserID(r.FormValue("user_id"))
		v, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.ctrl.PutRating(r.Context(), recordID, recordType, &model.Rating{UserID: userID, Value: model.RatingValue(v)}); err != nil {
			log.Printf("Repository put error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
