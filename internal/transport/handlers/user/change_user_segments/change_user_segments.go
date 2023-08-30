package change_user_segments

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
	"github.com/kiryu-dev/segments-api/internal/repository"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers"
	"github.com/kiryu-dev/segments-api/internal/transport/validation"
)

type segmentChanger interface {
	Change(context.Context, []*model.UserSegment, model.OpType) []error
}

type segmentWithTTL struct {
	Slug string  `json:"slug"`
	TTL  *string `json:"ttl"`
}

type segments []*segmentWithTTL

type request struct {
	UserID   uint64   `json:"user_id"`
	ToAdd    segments `json:"to_add"`
	ToDelete []string `json:"to_delete"`
}

type response struct {
	Slug       string `json:"slug"`
	OpType     string `json:"operation_type"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// ChangeUserSegments godoc
//
//	@Summary		Изменить сегменты пользователя
//	@Description	Метод изменения активных сегментов пользователя. Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя. Также есть возможность задать TTL для добавляемых сегментов, чтобы по истечению времени они автоматически удалились у пользователя. TTL задается в формате "1y8m21d". Если хотите только удалить определенные сегменты, то можно опустить список для добавления и наоборот.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			input	body		request					true	"user id, segment's list to add (with ttl optional), segment's list to delete"
//	@Success		200		{object}	response				"list of changes"
//	@Failure		400		{object}	handlers.responseError	"error"
//	@Failure		500		{object}	handlers.responseError	"error"
//	@Failure		default	{object}	handlers.responseError	"error"
//	@Router			/user-segments [post]
func New(service segmentChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var (
			data = new(request)
			err  = json.NewDecoder(r.Body).Decode(data)
		)
		defer r.Body.Close()
		if err != nil || len(data.ToAdd) == 0 && len(data.ToDelete) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, "invalid data to change user's segments")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		addSeg, err := data.ToAdd.toSegmentModel(data.UserID)
		if errors.Is(err, validation.ErrRegexpErr) {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteServerError(w, http.StatusInternalServerError)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		var (
			addErr = service.Change(ctx, addSeg, model.AddOp)
			offset = len(addErr)
			delErr = service.Change(ctx, slugsToSegment(data.UserID, data.ToDelete), model.DeleteOp)
			resp   = make([]*response, offset+len(delErr))
		)
		for i, err := range addErr {
			resp[i] = createResponse(err, data.ToAdd[i].Slug, model.AddOp)
		}
		for i, err := range delErr {
			resp[offset+i] = createResponse(err, data.ToDelete[i], model.DeleteOp)
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteServerError(w, http.StatusInternalServerError)
		}
	}
}

func createResponse(err error, slug string, op model.OpType) *response {
	resp := &response{
		Slug:       slug,
		StatusCode: http.StatusOK,
		OpType:     op.String(),
	}
	if errors.Is(err, repository.ErrSegmentNotExists) ||
		errors.Is(err, repository.ErrHasSegment) {
		resp.StatusCode = http.StatusBadRequest
		resp.Message = err.Error()
	} else if err != nil {
		resp.StatusCode = http.StatusInternalServerError
	}
	return resp
}

func slugsToSegment(id uint64, slugs []string) []*model.UserSegment {
	result := make([]*model.UserSegment, len(slugs))
	for i, slug := range slugs {
		result[i] = &model.UserSegment{
			UserID: id,
			Slug:   slug,
		}
	}
	return result
}

func (s segments) toSegmentModel(userID uint64) ([]*model.UserSegment, error) {
	result := make([]*model.UserSegment, len(s))
	for i, seg := range s {
		result[i] = &model.UserSegment{
			UserID: userID,
			Slug:   seg.Slug,
		}
		if seg.TTL == nil {
			continue
		}
		ttl, err := validation.ValidateTTL(*seg.TTL)
		if err != nil {
			return nil, err
		}
		deleteTime := time.Now().AddDate(ttl.Years, ttl.Months, ttl.Days)
		result[i].DeleteTime = &deleteTime
	}
	return result, nil
}
