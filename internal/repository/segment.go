package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/segments-api/internal/model"
	"github.com/kiryu-dev/segments-api/pkg/util/parser"
)

var (
	ErrSegmentExists    = fmt.Errorf("specified segment already exists")
	ErrSegmentNotExists = fmt.Errorf("specified segment doesn't exist")
	ErrHasSegment       = fmt.Errorf("user already has specified segment")
	ErrUserNotExists    = fmt.Errorf("user with specified id doesn't exist")
	ErrNoUsers          = fmt.Errorf("there're no users with specified segment")
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
	query := `
INSERT INTO users_segments (user_id, slug, delete_time)
VALUES ($1, $2, $3);
	`
	err := s.findDublicate(ctx, seg.UserID, seg.Slug)
	if err != sql.ErrNoRows {
		return ErrHasSegment
	}
	_, err = s.db.ExecContext(ctx, query, seg.UserID, seg.Slug, seg.DeleteTime)
	return err
}

func (s *segmentRepository) DeleteFromUser(ctx context.Context, seg *model.UserSegment) error {
	query := `DELETE FROM users_segments WHERE user_id = $1 AND slug = $2;`
	res, err := s.db.ExecContext(ctx, query, seg.UserID, seg.Slug)
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
	var (
		query    = `SELECT slug FROM users_segments WHERE user_id = $1;`
		segments = make([]string, 0)
	)
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting segments of user with ID %d: %v", userID, err)
	}
	for rows.Next() {
		var segment string
		if err := rows.Scan(&segment); err != nil {
			return nil, fmt.Errorf("error getting segments of user with ID %d: %v", userID, err)
		}
		segments = append(segments, segment)
	}
	if len(segments) == 0 {
		return nil, ErrUserNotExists
	}
	return segments, nil
}

func (s *segmentRepository) DeleteByTTL(ctx context.Context) ([]*model.UserSegment, error) {
	var (
		query = `
DELETE FROM users_segments WHERE delete_time < NOW()
RETURNING (user_id, slug);
		`
		segments = make([]*model.UserSegment, 0)
	)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error deleting time expired segments: %v", err)
	}
	for rows.Next() {
		buf := make([]byte, 0)
		if err := rows.Scan(&buf); err != nil {
			return nil, fmt.Errorf("error getting deleted users' segments: %v", err)
		}
		id, slug := parser.ParseResponse(buf)
		segments = append(segments, &model.UserSegment{
			UserID: id,
			Slug:   slug,
		})
	}
	return segments, nil
}

func (s *segmentRepository) GetUsersBySegment(ctx context.Context, slug string) ([]uint64, error) {
	var (
		query = `SELECT user_id FROM users_segments WHERE slug = $1;`
		users = make([]uint64, 0)
	)
	rows, err := s.db.QueryContext(ctx, query, slug)
	if err != nil {
		return nil, fmt.Errorf("error getting users by specified segment %s: %v", slug, err)
	}
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error getting users by specified segment %s: %v", slug, err)
		}
		users = append(users, id)
	}
	if len(users) == 0 {
		return nil, ErrNoUsers
	}
	return users, nil
}

func (s *segmentRepository) findDublicate(ctx context.Context, userID uint64, slug string) error {
	query := `SELECT user_id FROM users_segments WHERE user_id = $1 AND slug = $2;`
	return s.db.QueryRowContext(ctx, query, userID, slug).Scan()
}
