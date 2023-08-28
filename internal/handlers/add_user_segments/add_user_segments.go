package add_user_segments

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type segmentsAdder interface {
	AddToUser(context.Context, []*model.UserSegment) <-chan *model.ErrSegmentInfo
}

type duration struct {
	time.Duration
}

type segment struct {
	Slug string    `json:"slug"`
	TTL  *duration `json:"ttl"`
}

type request struct {
	UserID   uint64    `json:"user_id"`
	Segments []segment `json:"segments"`
}

func New(service segmentsAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		data := new(request)
		if err := json.Unmarshal(body, data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid data for adding user to segments"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		segments := make([]*model.UserSegment, len(data.Segments))
		for i, seg := range data.Segments {
			segments[i] = &model.UserSegment{
				UserID:     data.UserID,
				Slug:       seg.Slug,
				DeleteTime: seg.TTL.toDeleteTime(),
			}
		}
		var (
			ch   = service.AddToUser(ctx, segments)
			resp = make([]*model.ErrSegmentInfo, 0)
		)
		for errInfo := range ch {
			resp = append(resp, errInfo)
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
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
