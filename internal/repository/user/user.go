package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/segments-api/internal/repository"
)

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) *repo {
	return &repo{db}
}

func (r *repo) Create(ctx context.Context, userID uint64) error {
	query := `INSERT INTO users (id) VALUES ($1);`
	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return repository.ErrUserExists
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, userID uint64) error {
	query := `DELETE FROM users WHERE id = $1;`
	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user with ID %d: %v", userID, err)
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return repository.ErrUserNotExists
	}
	return nil
}

func (r *repo) GetUserSegments(ctx context.Context, userID uint64) ([]string, error) {
	var (
		query    = `SELECT slug FROM users_segments WHERE user_id = $1;`
		segments = make([]string, 0)
	)
	rows, err := r.db.QueryContext(ctx, query, userID)
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
		return nil, repository.ErrUserNotExists
	}
	return segments, nil
}
