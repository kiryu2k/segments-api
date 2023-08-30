package create_segment

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kiryu-dev/segments-api/internal/transport/handlers"
	"github.com/kiryu-dev/segments-api/internal/transport/validation"
)

type segmentCreator interface {
	Create(context.Context, string, float64) ([]uint64, error)
}

type request struct {
	Slug       string  `json:"slug"`
	Percentage float64 `json:"percentage"`
}

type response struct {
	Slug    string   `json:"slug"`
	UsersID []uint64 `json:"users_id"`
}

// CreateSegment godoc
//
//	@Summary		Создать новый сегмент
//	@Description	Метод создания сегмента. Принимает slug (название) сегмента. Опционально можно указать процент пользователей, которые добавятся в этот сегмент автоматически.
//	@Tags			segment
//	@Accept			json
//	@Produce		json
//	@Param			input	body		request					true	"segment name and user percentage (optional)"
//	@Success		200		{object}	response				"(optional) segment name and added users"
//	@Failure		400		{object}	handlers.responseError	"error"
//	@Failure		500		{object}	handlers.responseError	"error"
//	@Failure		default	{object}	handlers.responseError	"error"
//	@Router			/segment [post]
func New(service segmentCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := new(request)
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, "invalid data for segment creation")
			return
		}
		defer r.Body.Close()
		err := validation.ValidateSlug(data.Slug)
		if errors.Is(err, validation.ErrInvalidChar) || errors.Is(err, validation.ErrInvalidSize) {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteServerError(w, http.StatusInternalServerError)
			return
		}
		if err := validation.ValidatePercentage(data.Percentage); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		resp := &response{Slug: data.Slug}
		resp.UsersID, err = service.Create(ctx, data.Slug, data.Percentage)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteServerError(w, http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteServerError(w, http.StatusInternalServerError)
		}
	}
}
