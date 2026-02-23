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
	"google.golang.org/grpc"
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
		log.Printf("Starting http server on %s", cfg.App.HttpSocketStr())
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting http server: %v", err)
		}
	}()
	go func() {
		lis, err := net.Listen("tcp", cfg.App.GrpcSocketStr())
		if err != nil {
			log.Fatalf("failed to listen grpc: %v", err)
		}
		log.Printf("Starting gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	grpcServer.GracefulStop()

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
