package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/config"
	"github.com/kiryu-dev/segments-api/internal/handlers/create_segment"
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
		log.Fatal(err)
	}
	db, err := postgres.New(&cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
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
		log.Fatal(err)
	}
}

func setupRoutes(service *service.SegmentService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/segment", create_segment.New(service)).Methods(http.MethodPost)
	router.HandleFunc("/segment", nil).Methods(http.MethodDelete)
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
