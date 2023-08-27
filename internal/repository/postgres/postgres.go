package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/kiryu-dev/segments-api/internal/config"
	_ "github.com/lib/pq"
)

func New(cfg *config.DB) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.String())
	if err != nil {
		return nil, fmt.Errorf("invalid connection to postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot get access to postgtes: %w", err)
	}
	if err := initDB(db, cfg.InitFilepath); err != nil {
		return nil, fmt.Errorf("cannot initialize tables: %w", err)
	}
	return db, nil
}

func initDB(db *sql.DB, initFilepath string) error {
	buf, err := os.ReadFile(initFilepath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(buf))
	return err
}
