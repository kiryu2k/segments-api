package user

import (
	"context"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type userRepository interface {
	Create(context.Context, uint64) error
	Delete(context.Context, uint64) error
	GetUserSegments(context.Context, uint64) ([]string, error)
}

type logsRepository interface {
	Write(context.Context, *model.UserLog) error
}

type Service struct {
	repo userRepository
	logs logsRepository
}

func New(repo userRepository, logs logsRepository) *Service {
	return &Service{repo, logs}
}

func (s *Service) Create(ctx context.Context, userID uint64) error {
	return s.repo.Create(ctx, userID)
}

func (s *Service) Delete(ctx context.Context, userID uint64) error {
	slugs, _ := s.repo.GetUserSegments(ctx, userID)
	if err := s.repo.Delete(ctx, userID); err != nil {
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
	return s.repo.GetUserSegments(ctx, userID)
}
