package service

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
	DeleteByTTL(context.Context) error
}

type changeFunc func(context.Context, *model.UserSegment) error

type SegmentService struct {
	repo segmentRepository
}

type segmentError struct {
	idx int
	err error
}

func New(repo segmentRepository) *SegmentService {
	return &SegmentService{repo}
}

func (s *SegmentService) Create(ctx context.Context, slug string) error {
	return s.repo.Create(ctx, slug)
}

func (s *SegmentService) Delete(ctx context.Context, slug string) error {
	return s.repo.Delete(ctx, slug)
}

func (s *SegmentService) Change(ctx context.Context, seg []*model.UserSegment, opType int) []error {
	var (
		fn   changeFunc
		errs = make([]error, len(seg))
	)
	switch opType {
	case model.AddOp:
		fn = s.repo.AddToUser
	case model.DeleteOp:
		fn = s.repo.DeleteFromUser
	}
	for segErr := range changeSegments(ctx, seg, fn) {
		errs[segErr.idx] = segErr.err
	}
	return errs
}

func changeSegments(ctx context.Context, seg []*model.UserSegment,
	fn changeFunc) <-chan *segmentError {
	var (
		wg  = &sync.WaitGroup{}
		out = make(chan *segmentError)
	)
	wg.Add(len(seg))
	for i, s := range seg {
		go func(ctx context.Context, i int, seg *model.UserSegment) {
			defer wg.Done()
			out <- &segmentError{
				idx: i,
				err: fn(ctx, seg),
			}
		}(ctx, i, s)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (s *SegmentService) GetUserSegments(ctx context.Context, userID uint64) ([]string, error) {
	return s.repo.GetUserSegments(ctx, userID)
}

func (s *SegmentService) DeleteByTTL() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.repo.DeleteByTTL(ctx)
}
