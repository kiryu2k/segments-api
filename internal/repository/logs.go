package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type logsRepository struct {
	db *sql.DB
}

func NewLogger(db *sql.DB) *logsRepository {
	return &logsRepository{db}
}

func (l *logsRepository) Write(ctx context.Context, log *model.UserLog) error {
	query := `INSERT INTO logs (user_id, slug, operation, request_time) VALUES ($1, $2, $3, $4);`
	_, err := l.db.ExecContext(ctx, query, log.UserID, log.Slug, log.Operation, log.RequestTime)
	if err != nil {
		return fmt.Errorf("failed to write log of user %d with segment %s: %v", log.UserID, log.Slug, err)
	}
	return nil
}
