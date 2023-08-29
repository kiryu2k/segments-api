package logs

import (
	"context"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
	"github.com/kiryu-dev/segments-api/pkg/util/csv"
)

type logsRepository interface {
	Read(context.Context, uint64, time.Time) ([]*model.UserLog, error)
}

type Service struct {
	repo logsRepository
}

func New(repo logsRepository) *Service {
	return &Service{repo}
}

func (s *Service) GetUserLogs(ctx context.Context, userID uint64, date time.Time) (string, error) {
	logs, err := s.repo.Read(ctx, userID, date)
	if err != nil {
		return "", err
	}
	path, err := csv.GenerateCSV[model.UserLog](logs)
	if err != nil {
		return "", err
	}
	return path, nil
}
