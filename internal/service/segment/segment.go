package segment

import (
	"context"
	"sync"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type segmentRepository interface {
	Create(context.Context, string) error
	Delete(context.Context, string) error
	AddToUser(context.Context, *model.UserSegment) error
	DeleteFromUser(context.Context, *model.UserSegment) error
	GetUserSegments(context.Context, uint64) ([]string, error)
	DeleteByTTL(context.Context) ([]*model.UserSegment, error)
	GetUsersBySegment(context.Context, string) ([]uint64, error)
}

type logsRepository interface {
	Write(context.Context, *model.UserLog) error
}

type changeFunc func(context.Context, *model.UserSegment) error

type Service struct {
	repo segmentRepository
	logs logsRepository
}

type segmentError struct {
	idx int
	err error
}

func New(repo segmentRepository, logs logsRepository) *Service {
	return &Service{repo, logs}
}

func (s *Service) Create(ctx context.Context, slug string) error {
	return s.repo.Create(ctx, slug)
}

func (s *Service) Delete(ctx context.Context, slug string) error {
	users, _ := s.repo.GetUsersBySegment(ctx, slug)
	if err := s.repo.Delete(ctx, slug); err != nil {
		return err
	}
	for _, id := range users {
		_ = s.logs.Write(ctx, &model.UserLog{
			UserID:      id,
			Slug:        slug,
			Operation:   model.DeleteOp.String(),
			RequestTime: time.Now(),
		})
	}
	return nil
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
		fn = s.repo.AddToUser
	case model.DeleteOp:
		fn = s.repo.DeleteFromUser
	}
	return fn
}

func (s *Service) GetUserSegments(ctx context.Context,
	userID uint64) ([]string, error) {
	return s.repo.GetUserSegments(ctx, userID)
}

func (s *Service) DeleteByTTL() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	segments, err := s.repo.DeleteByTTL(ctx)
	if err != nil {
		return err
	}
	for _, segment := range segments {
		_ = s.logs.Write(ctx, &model.UserLog{
			UserID:      segment.UserID,
			Slug:        segment.Slug,
			Operation:   model.DeleteOp.String(),
			RequestTime: time.Now(),
		})
	}
	return nil
}
