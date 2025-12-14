package handler

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/middleware"
	"github.com/CXTACLYSM/hiring-api/internal/routing"
	"github.com/CXTACLYSM/hiring-api/pkg/postgres"
)

type LoginHandler struct {
	middlewares []middleware.Middleware
	connector   *postgres.Connector
}

func NewLoginHandler(connector *postgres.Connector, middlewares []middleware.Middleware) *LoginHandler {
	return &LoginHandler{
		connector:   connector,
		middlewares: middlewares,
	}
}

func (h *LoginHandler) Middlewares() []middleware.Middleware {
	return h.middlewares
}

func (h *LoginHandler) RouteMetadata() *routing.Metadata {
	return routing.NewRouteMetadata(routing.ApiPrefix, routing.V1, routing.RouteLogin)
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
