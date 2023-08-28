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

func (s *segmentRepository) AddToUser(ctx context.Context, seg *model.UserSegments) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var (
		searchQuery = `SELECT id FROM segment WHERE slug = $1`
		insertQuery = `
INSERT INTO users_segments (user_id, segment_id, delete_time)
VALUES ($1, $2, $3);
		`
	)
	for _, segment := range seg.Segments {
		var id uint64
		err := tx.QueryRowContext(ctx, searchQuery, segment.Slug).Scan(&id)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				return fmt.Errorf("transaction error: %w, rollback error: %w", err, rErr)
			}
			return ErrSegmentNotExists
		}
		_, err = tx.ExecContext(ctx, insertQuery, seg.UserID, id, segment.DeleteTime)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				return fmt.Errorf("transaction error: %w, rollback error: %w", err, rErr)
			}
			return err
		}
	}
	return tx.Commit()
}

func (s *segmentRepository) DeleteFromUser(ctx context.Context, seg *model.UserSegments) error {
	return nil
}

func (s *segmentRepository) GetActiveUserSegments(ctx context.Context, userID uint64) ([]string, error) {
	return nil, nil
}
