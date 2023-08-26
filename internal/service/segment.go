package service

import (
	"context"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type segmentRepository interface {
	Create(context.Context, string) error
	Delete(context.Context, string) error
	AddToUser(context.Context, *model.UserSegments) error
	DeleteFromUser(context.Context, *model.UserSegments) error
	GetActiveUserSegments(context.Context, uint64) ([]string, error)
}

type segmentService struct {
	repo segmentRepository
}

func New(repo segmentRepository) *segmentService {
	return &segmentService{repo}
}

func (s *segmentService) Create(ctx context.Context, slug string) error {
	return s.repo.Create(ctx, slug)
}

func (s *segmentService) Delete(ctx context.Context, slug string) error {
	return s.repo.Delete(ctx, slug)
}

func (s *segmentService) AddToUser(ctx context.Context, seg *model.UserSegments) error {
	if len(seg.Slugs) == 0 {
		return nil
	}
	return s.repo.AddToUser(ctx, seg)
}

func (s *segmentService) DeleteFromUser(ctx context.Context, seg *model.UserSegments) error {
	if len(seg.Slugs) == 0 {
		return nil
	}
	return s.repo.DeleteFromUser(ctx, seg)
}

func (s *segmentService) GetActiveUserSegments(ctx context.Context, userID uint64) ([]string, error) {
	return s.repo.GetActiveUserSegments(ctx, userID)
}
