package blog

import (
	"time"

	_ "github.com/CXTACLYSM/hiring-api/cmd/blog/docs"
	"github.com/CXTACLYSM/hiring-api/internal/blog/di"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouter(middlewares *di.Middlewares, handlers *di.Handlers) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middlewares.Json.ContentType)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/info", handlers.Info.ServeHTTP)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.Authenticate.Authenticate)

			r.Get("/posts", handlers.FindListPost.ServeHTTP)
			r.Post("/posts", handlers.CreateOnePost.ServeHTTP)
			r.Get("/posts/{id}", handlers.FindOnePost.ServeHTTP)
			r.Put("/posts/{id}", handlers.UpdateOnePost.ServeHTTP)
			r.Delete("/posts/{id}", handlers.DeleteOnePost.ServeHTTP)
		})
	})

	return r
}
