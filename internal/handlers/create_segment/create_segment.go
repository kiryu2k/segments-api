package create_segment

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type segmentCreator interface {
	Create(context.Context, string) error
}

type request struct {
	Slug string `json:"slug"`
}

func New(service segmentCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := new(request)
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
		}
	}
}
