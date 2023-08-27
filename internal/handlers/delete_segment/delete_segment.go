package delete_segment

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
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
		ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
		defer cancel()
		if err := service.Delete(ctx, data.Slug); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("deleting error"))
			return
		}
	}
}
