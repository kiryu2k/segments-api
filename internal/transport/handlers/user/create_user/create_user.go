package create_user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type userCreator interface {
	Create(context.Context, uint64) error
}

type request struct {
	UserID uint64 `json:"user_id"`
}

func New(service userCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := new(request)
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid data for user creation"))
			return
		}
		defer r.Body.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := service.Create(ctx, data.UserID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
