package delete_segment

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/repository"
	"github.com/kiryu-dev/segments-api/internal/validation"
)

type segmentDeleter interface {
	Delete(context.Context, string) error
}

func New(service segmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]
		if err := validation.ValidateSlug(slug); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		err := service.Delete(ctx, slug)
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
