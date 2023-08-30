package delete_segment

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/repository"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers"
	"github.com/kiryu-dev/segments-api/internal/transport/validation"
)

type segmentDeleter interface {
	Delete(context.Context, string) error
}

// DeleteSegment godoc
//
//	@Summary		Удалить сегмент
//	@Description	Метод удаления сегмента. Принимает slug (название) сегмента.
//	@Tags			segment
//	@Produce		json
//	@Param			slug	path	string	true	"segment name"
//	@Success		200
//	@Failure		400		{object}	handlers.responseError	"error"
//	@Failure		500		{object}	handlers.responseError	"error"
//	@Failure		default	{object}	handlers.responseError	"error"
//	@Router			/segment/{slug} [delete]
func New(service segmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		slug := mux.Vars(r)["slug"]
		err := validation.ValidateSlug(slug)
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
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		err = service.Delete(ctx, slug)
		if errors.Is(err, repository.ErrSegmentNotExists) {
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
