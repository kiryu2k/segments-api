package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/config"
	"github.com/kiryu-dev/segments-api/internal/handlers/create_segment"
	"github.com/kiryu-dev/segments-api/internal/handlers/delete_segment"
	"github.com/kiryu-dev/segments-api/internal/repository"
	"github.com/kiryu-dev/segments-api/internal/repository/postgres"
	"github.com/kiryu-dev/segments-api/internal/service"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/config.yaml", "config file path")
	flag.Parse()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	db, err := postgres.New(&cfg.DB)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer db.Close()
	var (
		repo    = repository.New(db)
		service = service.New(repo)
		router  = setupRoutes(service)
		server  = &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			WriteTimeout: cfg.Timeout,
			ReadTimeout:  cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		}
	)
	if err := server.ListenAndServe(); err != nil {
		logger.Error(err.Error())
	}
}

func setupRoutes(service *service.SegmentService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/segment", create_segment.New(service)).Methods(http.MethodPost)
	router.HandleFunc("/segment", delete_segment.New(service)).Methods(http.MethodDelete)
	router.HandleFunc("/user-segments", nil).Methods(http.MethodPost)
	router.HandleFunc("/user-segments", nil).Methods(http.MethodGet)
	return router
}

// type segmentService interface {
// 	Create(context.Context, string) error
// 	Delete(context.Context, string) error
// 	AddToUser(context.Context, *model.UserSegments) error
// 	DeleteFromUser(context.Context, *model.UserSegments) error
// 	GetActiveUserSegments(context.Context, uint64) ([]string, error)
// }
