package http

import (
	"encoding/json"
	"errors"
	"github.com/supersherm5/movie-app/metadata/internal/controller/metadata"
	"log"
	"net/http"

	"github.com/supersherm5/movie-app/metadata/internal/repository"
)

// Handler defines a movie metadata HTTP handler.
type Handler struct {
	ctrl *metadata.Controller
}

// New creates a new movie metadata HTTP handler.
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl}
}

// GetMetadata handles Get /metadata requests.
func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	m, err := h.ctrl.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrMetadataNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Response JSON encode error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
