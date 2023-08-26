package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kiryu-dev/segments-api/internal/config"
	"github.com/kiryu-dev/segments-api/internal/repository"
	"github.com/kiryu-dev/segments-api/internal/repository/postgres"
	"github.com/kiryu-dev/segments-api/internal/service"
	"github.com/kiryu-dev/segments-api/internal/transport"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/config.yaml", "config file path")
	flag.Parse()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	db, err := postgres.New(&cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	var (
		repo    = repository.New(db)
		service = service.New(repo)
		router  = transport.New(service)
		server  = &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			WriteTimeout: cfg.Timeout,
			ReadTimeout:  cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		}
	)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
