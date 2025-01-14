package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"movieexample.com/movie/internal/controller/movie"
	"net/http"
)

// Handler defines a movie service HTTP handler.
type Handler struct {
	ctrl *movie.Controller
}

// New creates a new movie service HTTP handler.
func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMovieDetails handles GET /movie requests.
func (h *Handler) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	details, err := h.ctrl.Get(r.Context(), id)
	if err != nil && errors.Is(err, movie.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("GetMovieDetails", slog.String("id", id), slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(details)
	if err != nil {
		slog.Error("GetMovieDetails", slog.String("id", id), slog.String("err", err.Error()))
	}
}
