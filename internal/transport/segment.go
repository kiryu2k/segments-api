package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/model"
)

type segmentService interface {
	Create(context.Context, string) error
	Delete(context.Context, string) error
	AddToUser(context.Context, *model.UserSegments) error
	DeleteFromUser(context.Context, *model.UserSegments) error
	GetActiveUserSegments(context.Context, uint64) ([]string, error)
}

type segmentRequest struct {
	Slug string `json:"slug"`
}

func New(service segmentService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/segment", createSegment(service)).Methods(http.MethodPost)
	router.HandleFunc("/segment", nil).Methods(http.MethodDelete)
	router.HandleFunc("/user-segments", nil).Methods(http.MethodPost)
	router.HandleFunc("/user-segments", nil).Methods(http.MethodGet)
	return router
}

func createSegment(service segmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := new(segmentRequest)
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid data for segment creation"))
			return
		}
		defer r.Body.Close()
		ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
		defer cancel()
		if err := service.Create(ctx, data.Slug); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("creating error"))
			return
		}
	}
}

func deleteSegment(service segmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := new(segmentRequest)
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid data for segment deletion"))
			return
		}
		defer r.Body.Close()
		ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
		defer cancel()
		if err := service.Delete(ctx, data.Slug); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("deleting error"))
			return
		}
	}
}
