package change_user_segments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
	"github.com/kiryu-dev/segments-api/internal/repository"
)

type segmentChanger interface {
	Change(context.Context, []*model.UserSegment, int) []error
}

type duration struct {
	time.Duration
}

type segmentWithTTL struct {
	Slug string    `json:"slug"`
	TTL  *duration `json:"ttl"`
}

type segments []*segmentWithTTL

type request struct {
	UserID   uint64   `json:"user_id"`
	ToAdd    segments `json:"to_add"`
	ToDelete []string `json:"to_delete"`
}

type response struct {
	Slug       string `json:"slug"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func New(service segmentChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		data := new(request)
		err = json.Unmarshal(body, data)
		if err != nil || len(data.ToAdd) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid data to change user's segments"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		var (
			errs = service.Change(ctx, data.ToAdd.makeSegmentModel(data.UserID), model.AddOp)
			resp = make([]*response, len(errs))
		)
		for i, err := range errs {
			resp[i] = &response{
				Slug:       data.ToAdd[i].Slug,
				StatusCode: http.StatusOK,
			}
			if errors.Is(err, repository.ErrSegmentNotExists) ||
				errors.Is(err, repository.ErrHasSegment) {
				resp[i].StatusCode = http.StatusBadRequest
				resp[i].Message = err.Error()
				continue
			}
			if err != nil {
				resp[i].StatusCode = http.StatusInternalServerError
			}

		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (s segments) makeSegmentModel(userID uint64) []*model.UserSegment {
	segments := make([]*model.UserSegment, len(s))
	for i, seg := range s {
		segments[i] = &model.UserSegment{
			UserID:     userID,
			Slug:       seg.Slug,
			DeleteTime: seg.TTL.toDeleteTime(),
		}
	}
	return segments
}

func (d *duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		return err
	}
	return fmt.Errorf("invalid duration")
}

func (d *duration) toDeleteTime() *time.Time {
	if d != nil {
		deleteTime := time.Now().Add(d.Duration)
		return &deleteTime
	}
	return nil
}
