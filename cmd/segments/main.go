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
	logs_repo "github.com/kiryu-dev/segments-api/internal/repository/logs"
	"github.com/kiryu-dev/segments-api/internal/repository/postgres"
	segment_repo "github.com/kiryu-dev/segments-api/internal/repository/segment"
	user_repo "github.com/kiryu-dev/segments-api/internal/repository/user"
	"github.com/kiryu-dev/segments-api/internal/service/logs"
	logs_service "github.com/kiryu-dev/segments-api/internal/service/logs"
	"github.com/kiryu-dev/segments-api/internal/service/segment"
	segment_service "github.com/kiryu-dev/segments-api/internal/service/segment"
	user_service "github.com/kiryu-dev/segments-api/internal/service/user"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers/logs/get_user_logs"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers/segment/create_segment"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers/segment/delete_segment"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers/user/change_user_segments"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers/user/create_user"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers/user/delete_user"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers/user/get_user_segments"
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
		/* repository layer */
		logRepo     = logs_repo.New(db)
		userRepo    = user_repo.New(db)
		segmentRepo = segment_repo.New(db)
		/* service layer */
		logService     = logs_service.New(logRepo)
		userService    = user_service.New(userRepo, logRepo)
		segmentService = segment_service.New(segmentRepo, userService, logRepo)
		/* transport layer */
		router = setupRoutes(segmentService, userService, logService)
		server = &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			WriteTimeout: cfg.Timeout,
			ReadTimeout:  cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		}
	)
	go func() {
		for {
			if err := segmentService.DeleteByTTL(); err != nil {
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

func setupRoutes(segment *segment.Service, user *user_service.Service, log *logs.Service) *mux.Router {
	router := mux.NewRouter()
	{
		router.HandleFunc("/segment", create_segment.New(segment)).Methods(http.MethodPost)
		router.HandleFunc("/segment/{slug}", delete_segment.New(segment)).Methods(http.MethodDelete)
	}
	{
		router.HandleFunc("/user", create_user.New(user)).Methods(http.MethodPost)
		router.HandleFunc("/user/{userID}", delete_user.New(user)).Methods(http.MethodDelete)
		router.HandleFunc("/user-segments", change_user_segments.New(user)).Methods(http.MethodPost)
		router.HandleFunc("/user-segments/{userID}", get_user_segments.New(user)).Methods(http.MethodGet)
	}
	{
		router.HandleFunc("/log/{userID}", get_user_logs.New(log)).Methods(http.MethodGet)
	}
	return router
}
