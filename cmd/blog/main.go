package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CXTACLYSM/hiring-api/configs/blog"
	"github.com/CXTACLYSM/hiring-api/internal/blog"
	"github.com/CXTACLYSM/hiring-api/internal/blog/di"
)

func main() {
	cfg, err := configs.Create()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	container := &di.Container{}
	if err = container.Init(cfg); err != nil {
		log.Fatalf("Error initializing container: %s", err.Error())
	}
	defer container.Infrastructure.PgConnector.Close()
	defer container.Infrastructure.RedisConnector.Close()

	r := blog.InitRouter(container.Middlewares, container.Handlers)
	srv := http.Server{
		Addr:              cfg.App.HttpSocketStr(),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      35 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	go func() {
		log.Printf("Starting http server on %s", cfg.App.HttpSocketStr())
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting http server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
