package delete_segment

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kiryu-dev/segments-api/internal/repository"
)

type segmentDeleter interface {
	Delete(context.Context, string) error
}

type request struct {
	Slug string `json:"slug"`
}

func New(service segmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := new(request)
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid data for segment deletion"))
			return
		}
		defer r.Body.Close()
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		err := service.Delete(ctx, data.Slug)
		if errors.Is(err, repository.ErrSegmentNotExists) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("requested segment doesn't exist"))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
