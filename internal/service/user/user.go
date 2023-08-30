package user

import (
	"context"
	"sync"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type userRepository interface {
	Create(context.Context, uint64) error
	Delete(context.Context, uint64) error
	GetUserSegments(context.Context, uint64) ([]string, error)
	AddSegment(context.Context, *model.UserSegment) error
	DeleteSegment(context.Context, *model.UserSegment) error
}

type logsRepository interface {
	Write(context.Context, *model.UserLog) error
}

type Service struct {
	user userRepository
	logs logsRepository
}

type segmentError struct {
	idx int
	err error
}

type changeFunc func(context.Context, *model.UserSegment) error

func New(user userRepository, logs logsRepository) *Service {
	return &Service{user, logs}
}

func (s *Service) Create(ctx context.Context, userID uint64) error {
	return s.user.Create(ctx, userID)
}

func (s *Service) Delete(ctx context.Context, userID uint64) error {
	slugs, _ := s.user.GetUserSegments(ctx, userID)
	if err := s.user.Delete(ctx, userID); err != nil {
		return err
	}
	for _, slug := range slugs {
		_ = s.logs.Write(ctx, &model.UserLog{
			UserID:      userID,
			Slug:        slug,
			Operation:   model.DeleteOp.String(),
			RequestTime: time.Now(),
		})
	}
	return nil
}

func (s *Service) GetUserSegments(ctx context.Context, userID uint64) ([]string, error) {
	return s.user.GetUserSegments(ctx, userID)
}

func (s *Service) Change(ctx context.Context, seg []*model.UserSegment, opType model.OpType) []error {
	var (
		result  = make([]error, len(seg))
		errChan = s.changeSegments(ctx, seg, opType)
	)
	for e := range errChan {
		result[e.idx] = e.err
	}
	return result
}

func (s *Service) changeSegments(ctx context.Context, seg []*model.UserSegment,
	opType model.OpType) <-chan *segmentError {
	var (
		fn  = s.defineChangeFunc(opType)
		wg  = &sync.WaitGroup{}
		out = make(chan *segmentError)
	)
	wg.Add(len(seg))
	operation := opType.String()
	for i, segment := range seg {
		go func(ctx context.Context, i int, segment *model.UserSegment) {
			defer wg.Done()
			err := fn(ctx, segment)
			if err == nil {
				_ = s.logs.Write(ctx, &model.UserLog{
					UserID:      segment.UserID,
					Slug:        segment.Slug,
					Operation:   operation,
					RequestTime: time.Now(),
				})
			}
			out <- &segmentError{
				idx: i,
				err: err,
			}
		}(ctx, i, segment)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (s *Service) defineChangeFunc(opType model.OpType) changeFunc {
	var fn changeFunc
	switch opType {
	case model.AddOp:
		fn = s.user.AddSegment
	case model.DeleteOp:
		fn = s.user.DeleteSegment
	}
	return fn
}
