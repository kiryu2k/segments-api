package service

import (
	"context"
	"sync"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type segmentRepository interface {
	Create(context.Context, string) error
	Delete(context.Context, string) error
	AddToUser(context.Context, *model.UserSegment) error
	DeleteFromUser(context.Context, *model.UserSegment) error
	GetActiveUserSegments(context.Context, uint64) ([]string, error)
}

type SegmentService struct {
	repo segmentRepository
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

func (s *SegmentService) AddToUser(ctx context.Context, segments []*model.UserSegment) <-chan *model.ErrSegmentInfo {
	var (
		wg  = &sync.WaitGroup{}
		out = make(chan *model.ErrSegmentInfo)
	)
	wg.Add(len(segments))
	for _, seg := range segments {
		go func(ctx context.Context, seg *model.UserSegment) {
			defer wg.Done()
			err := s.repo.AddToUser(ctx, seg)
			if err != nil {
				out <- &model.ErrSegmentInfo{
					Slug:    seg.Slug,
					Message: err.Error(),
				}
			}
		}(ctx, seg)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (s *SegmentService) DeleteFromUser(ctx context.Context, seg *model.UserSegment) error {
	// if len(seg.Slugs) == 0 {
	// 	return nil
	// }
	return s.repo.DeleteFromUser(ctx, seg)
}

func (s *SegmentService) GetActiveUserSegments(ctx context.Context, userID uint64) ([]string, error) {
	return s.repo.GetActiveUserSegments(ctx, userID)
}
