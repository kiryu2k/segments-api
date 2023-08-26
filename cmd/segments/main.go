package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kiryu-dev/segments-api/internal/config"
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
	var (
		service = service.New(nil)
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
