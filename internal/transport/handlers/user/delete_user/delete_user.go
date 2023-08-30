package delete_user

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/repository"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers"
)

type userDeleter interface {
	Delete(context.Context, uint64) error
}

// DeleteUser godoc
//
//	@Summary		Удалить пользователя
//	@Description	Метод удаления пользователя. Принимает на вход id пользователя.
//	@Tags			user
//	@Produce		json
//	@Param			userID	path	int	true	"user id"
//	@Success		200
//	@Failure		400		{object}	handlers.responseError	"error"
//	@Failure		500		{object}	handlers.responseError	"error"
//	@Failure		default	{object}	handlers.responseError	"error"
//	@Router			/user/{userID} [delete]
func New(service userDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idStr := mux.Vars(r)["userID"]
		userID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		err = service.Delete(ctx, userID)
		if errors.Is(err, repository.ErrUserNotExists) {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteServerError(w, http.StatusInternalServerError)
		}
	}
}
