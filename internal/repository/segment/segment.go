package segment

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/segments-api/internal/model"
	"github.com/kiryu-dev/segments-api/internal/repository"
	"github.com/kiryu-dev/segments-api/pkg/util/parser"
)

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) *repo {
	return &repo{db}
}

func (r *repo) Create(ctx context.Context, slug string) error {
	query := `INSERT INTO segment (slug) VALUES ($1);`
	if _, err := r.db.ExecContext(ctx, query, slug); err != nil {
		return repository.ErrSegmentExists
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, slug string) error {
	query := `DELETE FROM segment WHERE slug = $1;`
	res, err := r.db.ExecContext(ctx, query, slug)
	if err != nil {
		return fmt.Errorf("error deleting segment with name %s: %v", slug, err)
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return repository.ErrSegmentNotExists
	}
	return nil
}

func (r *repo) DeleteByTTL(ctx context.Context) ([]*model.UserSegment, error) {
	var (
		query = `
DELETE FROM users_segments WHERE delete_time < NOW()
RETURNING (user_id, slug);
		`
		segments = make([]*model.UserSegment, 0)
	)
	rows, err := r.db.QueryContext(ctx, query)
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

func (r *repo) GetUsersBySegment(ctx context.Context, slug string) ([]uint64, error) {
	var (
		query = `SELECT user_id FROM users_segments WHERE slug = $1;`
		users = make([]uint64, 0)
	)
	rows, err := r.db.QueryContext(ctx, query, slug)
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
		return nil, repository.ErrNoUsers
	}
	return users, nil
}
