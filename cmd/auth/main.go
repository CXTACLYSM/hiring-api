package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CXTACLYSM/hiring-api/configs"
	"github.com/CXTACLYSM/hiring-api/internal/auth/di"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := configs.Create()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	container := di.NewContainer()
	err = container.Init(cfg)
	if err != nil {
		log.Fatalf("Error initializing container: %s", err.Error())
	}
	defer container.Infrastructure.PgConnector.Close()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(container.Middlewares.Json.ContentType)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", container.Handlers.Register.ServeHTTP)
		r.Post("/login", container.Handlers.Login.ServeHTTP)

		r.Group(func(r chi.Router) {
			r.Use(container.Middlewares.Authenticate.Authenticate)
			r.Get("/me", container.Handlers.Me.ServeHTTP)
		})
	})

	fmt.Printf("Starting http server on %s\n", cfg.App.SocketStr())
	err = http.ListenAndServe(cfg.App.SocketStr(), r)
	if err != nil {
		log.Fatalf("error creating http server: %v", err)
	}
}
