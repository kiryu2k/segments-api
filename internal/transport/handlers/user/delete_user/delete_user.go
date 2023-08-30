package delete_user

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/repository"
)

type userDeleter interface {
	Delete(context.Context, uint64) error
}

func New(service userDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["userID"]
		userID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid user id"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		err = service.Delete(ctx, userID)
		if errors.Is(err, repository.ErrUserNotExists) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
