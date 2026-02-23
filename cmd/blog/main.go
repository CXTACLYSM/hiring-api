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
	"go.uber.org/zap"
)

// @title           Blog API
// @version         1.0
// @description     Blog post management service.

// @host            localhost:8081
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}"

func main() {
	cfg, err := configs.Create()
	if err != nil {
		log.Fatalf("error creating config: %v", err)
	}

	container := &di.Container{}
	if err = container.Init(cfg); err != nil {
		log.Fatalf("error initializing container: %v", err)
	}

	logger := container.Infrastructure.Logger
	defer logger.Sync()
	defer container.Infrastructure.PgConnector.Close()
	defer container.Infrastructure.RedisConnector.Close()
	defer container.Infrastructure.AuthGrpcConn.Close()

	r := blog.InitRouter(container.Middlewares, container.Handlers)
	srv := http.Server{
		Addr:              cfg.App.HttpSocketStr(),
		Handler:           r,
		ReadHeaderTimeout: cfg.App.Http.ReadHeaderTimeout,
		ReadTimeout:       cfg.App.Http.ReadTimeout,
		WriteTimeout:      cfg.App.Http.WriteTimeout,
		IdleTimeout:       cfg.App.Http.IdleTimeout,
		MaxHeaderBytes:    cfg.App.Http.MaxHeaderBytes,
	}

	go func() {
		logger.Info("starting http server on %s", zap.String("addr", cfg.App.HttpSocketStr()))
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("error starting http server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server stopped")
}
