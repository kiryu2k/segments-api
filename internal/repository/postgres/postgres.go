package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/kiryu-dev/segments-api/internal/config"
	_ "github.com/lib/pq"
)

func New(cfg *config.DB) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.String())
	if err != nil {
		return nil, fmt.Errorf("invalid connection to postgres: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("cannot get access to postgtes: %w", err)
	}
	if err := initDB(ctx, db, cfg.InitFilepath); err != nil {
		return nil, fmt.Errorf("cannot initialize tables: %w", err)
	}
	return db, nil
}

func initDB(ctx context.Context, db *sql.DB, initFilepath string) error {
	buf, err := os.ReadFile(initFilepath)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, string(buf))
	return err
}
