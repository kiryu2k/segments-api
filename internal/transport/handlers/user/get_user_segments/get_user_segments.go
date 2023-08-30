package get_user_segments

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers"
)

type segmentsGetter interface {
	GetUserSegments(context.Context, uint64) ([]string, error)
}

// GetUserSegments godoc
//
//	@Summary		Получить активные сегменты пользователя
//	@Description	Метод получения активных сегментов пользователя. Принимает на вход id пользователя.
//	@Tags			user
//	@Produce		json
//	@Param			userID	path		int						true	"user id"
//	@Success		200		{array}		string					"list of segments"
//	@Failure		400		{object}	handlers.responseError	"error"
//	@Failure		500		{object}	handlers.responseError	"error"
//	@Failure		default	{object}	handlers.responseError	"error"
//	@Router			/user-segments/{userID} [get]
func New(service segmentsGetter) http.HandlerFunc {
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
		segments, err := service.GetUserSegments(ctx, userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(segments); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteServerError(w, http.StatusInternalServerError)
		}
	}
}
