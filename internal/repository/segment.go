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
	ErrHasSegment       = fmt.Errorf("user already has specified segment")
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
		query = `
INSERT INTO users_segments (user_id, segment_id, delete_time)
VALUES ($1, $2, $3);
		`
		segmentID, err = s.getSegmentID(ctx, seg.Slug)
	)
	if err == sql.ErrNoRows {
		return ErrSegmentNotExists
	}
	if err != nil {
		return fmt.Errorf("error adding segment %s to user with ID %d: %v",
			seg.Slug, seg.UserID, err)
	}
	if err := s.findDublicate(ctx, seg.UserID, segmentID); err != sql.ErrNoRows {
		return ErrHasSegment
	}
	_, err = s.db.ExecContext(ctx, query, seg.UserID, segmentID, seg.DeleteTime)
	return err
}

func (s *segmentRepository) DeleteFromUser(ctx context.Context, seg *model.UserSegment) error {
	var (
		query = `
DELETE FROM users_segments
WHERE user_id = $1 AND segment_id = $2;
		`
		segmentID, err = s.getSegmentID(ctx, seg.Slug)
	)
	if err == sql.ErrNoRows {
		return ErrSegmentNotExists
	}
	if err != nil {
		return fmt.Errorf("error deleting segment %s to user with ID %d: %v",
			seg.Slug, seg.UserID, err)
	}
	res, err := s.db.ExecContext(ctx, query, seg.UserID, segmentID)
	if err != nil {
		return fmt.Errorf("error deleting segment %s to user with ID %d: %v",
			seg.Slug, seg.UserID, err)
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return ErrSegmentNotExists
	}
	return nil
}

func (s *segmentRepository) GetUserSegments(ctx context.Context, userID uint64) ([]string, error) {
	return nil, nil
}

func (s *segmentRepository) getSegmentID(ctx context.Context, slug string) (uint64, error) {
	var (
		id    uint64
		query = `SELECT id FROM segment WHERE slug = $1;`
	)
	return id, s.db.QueryRowContext(ctx, query, slug).Scan(&id)
}

func (s *segmentRepository) findDublicate(ctx context.Context, userID, segmentID uint64) error {
	query := `SELECT user_id FROM users_segments WHERE user_id = $1 AND segment_id = $2;`
	return s.db.QueryRowContext(ctx, query, userID, segmentID).Scan()
}
