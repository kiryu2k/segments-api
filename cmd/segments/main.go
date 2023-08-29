package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/config"
	"github.com/kiryu-dev/segments-api/internal/handlers/change_user_segments"
	"github.com/kiryu-dev/segments-api/internal/handlers/create_segment"
	"github.com/kiryu-dev/segments-api/internal/handlers/delete_segment"
	"github.com/kiryu-dev/segments-api/internal/handlers/get_user_segments"
	"github.com/kiryu-dev/segments-api/internal/repository"
	"github.com/kiryu-dev/segments-api/internal/repository/postgres"
	"github.com/kiryu-dev/segments-api/internal/service"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/config.yaml", "config file path")
	flag.Parse()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("cannot load app's configuration: %v", err)
		return
	}
	log.Println("connecting to database...")
	db, err := postgres.New(&cfg.DB)
	if err != nil {
		log.Printf("unexpected database error: %v", err)
		return
	}
	defer db.Close()
	var (
		logger  = repository.NewLogger(db)
		repo    = repository.New(db)
		service = service.New(repo, logger)
		router  = setupRoutes(service)
		server  = &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			WriteTimeout: cfg.Timeout,
			ReadTimeout:  cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		}
	)
	go func() {
		for {
			if err := service.DeleteByTTL(); err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	go func() {
		log.Println("server is starting...")
		if err := server.ListenAndServe(); err != nil {
			log.Printf("failed to start server: %v", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("gracefully shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}
}

func setupRoutes(service *service.SegmentService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/segment", create_segment.New(service)).Methods(http.MethodPost)
	router.HandleFunc("/segment/{slug}", delete_segment.New(service)).Methods(http.MethodDelete)
	router.HandleFunc("/user-segments", change_user_segments.New(service)).Methods(http.MethodPost)
	router.HandleFunc("/user-segments/{userID}", get_user_segments.New(service)).Methods(http.MethodGet)
	return router
}
