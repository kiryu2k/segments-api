package logs

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
)

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) *repo {
	return &repo{db}
}

func (r *repo) Write(ctx context.Context, log *model.UserLog) error {
	query := `INSERT INTO logs (user_id, slug, operation, request_time) VALUES ($1, $2, $3, $4);`
	_, err := r.db.ExecContext(ctx, query, log.UserID, log.Slug, log.Operation, log.RequestTime)
	if err != nil {
		return fmt.Errorf("failed to write log of user %d with segment %s: %v", log.UserID, log.Slug, err)
	}
	return nil
}

func (r *repo) Read(ctx context.Context, userID uint64, date time.Time) ([]*model.UserLog, error) {
	var (
		query = `
SELECT * FROM logs WHERE user_id = $1
AND EXTRACT(YEAR FROM request_time) = $2
AND EXTRACT(MONTH FROM request_time) = $3
ORDER BY request_time;
		`
		logs = make([]*model.UserLog, 0)
	)
	rows, err := r.db.QueryContext(ctx, query, userID, date.Year(), int(date.Month()))
	if err != nil {
		return nil, fmt.Errorf("error getting logs of user %d: %v", userID, err)
	}
	for rows.Next() {
		log := new(model.UserLog)
		err := rows.Scan(&log.UserID, &log.Slug, &log.Operation, &log.RequestTime)
		if err != nil {
			return nil, fmt.Errorf("error getting logs of user %d: %v", userID, err)
		}
		logs = append(logs, log)
	}
	return logs, nil
}
