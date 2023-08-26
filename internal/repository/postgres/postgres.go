package postgres

import (
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/segments-api/internal/config"
)

func New(cfg *config.DB) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("invalid connection to postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot get access to postgtes: %w", err)
	}
	return db, nil
}
