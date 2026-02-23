package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CXTACLYSM/hiring-api/configs/auth"
	"github.com/CXTACLYSM/hiring-api/internal/auth"
	"github.com/CXTACLYSM/hiring-api/internal/auth/di"
	authgrpc "github.com/CXTACLYSM/hiring-api/internal/auth/grpc"
	pb "github.com/CXTACLYSM/hiring-api/pkg/grpc/auth/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// @title           Auth API
// @version         1.0
// @description     Authentication and user management service.

// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}" (with quotes around the token removed)

func main() {
	cfg, err := configs.Create()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	container := &di.Container{}
	if err = container.Init(cfg); err != nil {
		log.Fatalf("Error initializing container: %v", err)
	}

	logger := container.Infrastructure.Logger
	defer logger.Sync()
	defer container.Infrastructure.PgConnector.Close()
	defer container.Infrastructure.RedisConnector.Close()
	defer container.Infrastructure.Kafka.Publisher.Close()

	r := auth.InitRouter(container.Middlewares, container.Handlers)
	srv := &http.Server{
		Addr:              cfg.App.HttpSocketStr(),
		Handler:           r,
		ReadHeaderTimeout: cfg.App.Http.ReadHeaderTimeout,
		ReadTimeout:       cfg.App.Http.ReadTimeout,
		WriteTimeout:      cfg.App.Http.WriteTimeout,
		IdleTimeout:       cfg.App.Http.IdleTimeout,
		MaxHeaderBytes:    cfg.App.Http.MaxHeaderBytes,
	}
	grpcServer := grpc.NewServer()
	authServer := authgrpc.NewAuthServer(
		[]byte(cfg.App.JwtSecret),
		container.Queries.FindOneUser,
	)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	go func() {
		logger.Info("starting http server", zap.String("addr", cfg.App.HttpSocketStr()))
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server error", zap.Error(err))
		}
	}()
	go func() {
		lis, err := net.Listen("tcp", cfg.App.GrpcSocketStr())
		if err != nil {
			logger.Fatal("failed to listen grpc", zap.Error(err))
		}
		logger.Info("starting grpc server", zap.String("addr", cfg.App.GrpcSocketStr()))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve grpc", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	grpcServer.GracefulStop()

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server stopped")
}
