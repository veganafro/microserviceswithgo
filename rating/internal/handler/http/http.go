package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"movieexample.com/rating/internal/controller/rating"
	"movieexample.com/rating/pkg/model"
	"net/http"
	"strconv"
)

// Handler defines a rating service HTTP handler.
type Handler struct {
	ctrl *rating.Controller
}

// New creates a new rating service HTTP handler.
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

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
		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = json.NewEncoder(w).Encode(v)
		if err != nil {
			slog.Error("GetAggregatedRating error", slog.String("err", err.Error()))
			return
		}
	case http.MethodPut:
		userID := model.UserID(r.FormValue("userId"))
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		v, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		record := model.Rating{
			RecordID:   string(recordID),
			RecordType: string(recordType),
			UserID:     userID,
			Value:      model.RatingValue(v),
		}
		err = h.ctrl.PutRating(r.Context(), recordID, recordType, &record)
		if err != nil {
			slog.Error("PutRating error", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
