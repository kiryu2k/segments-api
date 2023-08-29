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
	Change(context.Context, []*model.UserSegment, model.OpType) []error
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
	OpType     string `json:"operation_type"`
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
		if err != nil || len(data.ToAdd) == 0 && len(data.ToDelete) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid data to change user's segments"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		var (
			addErr = service.Change(ctx, data.ToAdd.toSegmentModel(data.UserID), model.AddOp)
			offset = len(addErr)
			delErr = service.Change(ctx, slugsToSegment(data.UserID, data.ToDelete), model.DeleteOp)
			resp   = make([]*response, offset+len(delErr))
		)
		w.Header().Set("Content-Type", "encoding/json")
		for i, err := range addErr {
			resp[i] = createResponse(err, data.ToAdd[i].Slug, model.AddOp)
		}
		for i, err := range delErr {
			resp[offset+i] = createResponse(err, data.ToDelete[i], model.DeleteOp)
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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

func (s segments) toSegmentModel(userID uint64) []*model.UserSegment {
	result := make([]*model.UserSegment, len(s))
	for i, seg := range s {
		result[i] = &model.UserSegment{
			UserID:     userID,
			Slug:       seg.Slug,
			DeleteTime: seg.TTL.toDeleteTime(),
		}
	}
	return result
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
