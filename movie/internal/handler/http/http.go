package http

import (
	"encoding/json"
	"errors"
	"github.com/supersherm5/movie-app/movie/internal/controller/movie"
	"github.com/supersherm5/movie-app/movie/internal/gateway"
	"log"
	"net/http"
)

// Handler defines a movie handler.
type Handler struct {
	ctrl *movie.Controller
}

// New creates a new movie HTTP handler.
func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl}
}

// GetMovieDetails handles GET /movie request.
func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	details, err := h.ctrl.Get(req.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, movie.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			return
		case errors.Is(err, gateway.ErrServiceNotReachable):
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v", err)
	}
}
