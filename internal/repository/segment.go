package repository

import (
	"context"
	"database/sql"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type segmentRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *segmentRepository {
	return &segmentRepository{db}
}

func (s *segmentRepository) Create(ctx context.Context, slug string) error {
	return nil
}

func (s *segmentRepository) Delete(ctx context.Context, slug string) error {
	return nil
}

func (s *segmentRepository) AddToUser(ctx context.Context, seg *model.UserSegments) error {
	return nil
}

func (s *segmentRepository) DeleteFromUser(ctx context.Context, seg *model.UserSegments) error {
	return nil
}

func (s *segmentRepository) GetActiveUserSegments(ctx context.Context, userID uint64) ([]string, error) {
	return nil, nil
}
