package segment

import (
	"context"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type segmentRepository interface {
	Create(context.Context, string) error
	Delete(context.Context, string) error
	DeleteByTTL(context.Context) ([]*model.UserSegment, error)
	GetUsersBySegment(context.Context, string) ([]uint64, error)
}

type userRepository interface {
	GetUserSegments(context.Context, uint64) ([]string, error)
}

type logsRepository interface {
	Write(context.Context, *model.UserLog) error
}

type Service struct {
	segment segmentRepository
	user    userRepository
	logs    logsRepository
}

func New(segment segmentRepository, user userRepository, logs logsRepository) *Service {
	return &Service{segment, user, logs}
}

func (s *Service) Create(ctx context.Context, slug string) error {
	return s.segment.Create(ctx, slug)
}

func (s *Service) Delete(ctx context.Context, slug string) error {
	users, _ := s.segment.GetUsersBySegment(ctx, slug)
	if err := s.segment.Delete(ctx, slug); err != nil {
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

func (s *Service) DeleteByTTL() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	segments, err := s.segment.DeleteByTTL(ctx)
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
