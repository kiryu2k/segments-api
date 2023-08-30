package create_user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kiryu-dev/segments-api/internal/transport/handlers"
)

type userCreator interface {
	Create(context.Context, uint64) error
}

type request struct {
	UserID uint64 `json:"user_id"`
}

// CreateUser godoc
//
//	@Summary		Создать нового пользователя
//	@Description	Метод создания пользователя. Принимает на вход id пользователя.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			input	body	request	true	"user id"
//	@Success		200
//	@Failure		400		{object}	handlers.responseError	"error"
//	@Failure		500		{object}	handlers.responseError	"error"
//	@Failure		default	{object}	handlers.responseError	"error"
//	@Router			/user [post]
func New(service userCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := new(request)
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, "invalid data for user creation")
			return
		}
		defer r.Body.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := service.Create(ctx, data.UserID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteServerError(w, http.StatusInternalServerError)
		}
	}
}
