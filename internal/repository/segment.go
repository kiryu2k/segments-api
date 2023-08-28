package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/segments-api/internal/model"
)

var (
	ErrSegmentExists    = fmt.Errorf("specified segment already exists")
	ErrSegmentNotExists = fmt.Errorf("specified segment doesn't exist")
)

type segmentRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *segmentRepository {
	return &segmentRepository{db}
}

func (s *segmentRepository) Create(ctx context.Context, slug string) error {
	query := `INSERT INTO segment (slug) VALUES ($1);`
	if _, err := s.db.ExecContext(ctx, query, slug); err != nil {
		return ErrSegmentExists
	}
	return nil
}

func (s *segmentRepository) Delete(ctx context.Context, slug string) error {
	query := `DELETE FROM segment WHERE slug = $1;`
	res, err := s.db.ExecContext(ctx, query, slug)
	if err != nil {
		return fmt.Errorf("error deleting segment with name %s: %v", slug, err)
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return ErrSegmentNotExists
	}
	return nil
}

func (s *segmentRepository) AddToUser(ctx context.Context, seg *model.UserSegment) error {
	var (
		segmentID     uint64
		searchIDQuery = `SELECT id FROM segment WHERE slug = $1`
		insertQuery   = `
INSERT INTO users_segments (user_id, segment_id, delete_time)
VALUES ($1, $2, $3);
		`
	)
	err := s.db.QueryRowContext(ctx, searchIDQuery, seg.Slug).Scan(&segmentID)
	if err == sql.ErrNoRows {
		return ErrSegmentNotExists
	}
	if err != nil {
		return fmt.Errorf("error adding segment %s to user with ID %d: %v",
			seg.Slug, seg.UserID, err)
	}
	_, err = s.db.ExecContext(ctx, insertQuery, seg.UserID, segmentID, seg.DeleteTime)
	return err
}

func (s *segmentRepository) DeleteFromUser(ctx context.Context, seg *model.UserSegment) error {
	return nil
}

func (s *segmentRepository) GetActiveUserSegments(ctx context.Context, userID uint64) ([]string, error) {
	return nil, nil
}
