package handler

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/middleware"
	"github.com/CXTACLYSM/hiring-api/internal/routing"
	"github.com/CXTACLYSM/hiring-api/pkg/postgres"
)

type RegisterHandler struct {
	middlewares []middleware.Middleware
	connector   *postgres.Connector
}

func NewRegisterHandler(connector *postgres.Connector, middlewares []middleware.Middleware) *RegisterHandler {
	return &RegisterHandler{
		connector:   connector,
		middlewares: middlewares,
	}
}

func (h *RegisterHandler) Middlewares() []middleware.Middleware {
	return h.middlewares
}

func (h *RegisterHandler) RouteMetadata() *routing.Metadata {
	return routing.NewRouteMetadata(routing.ApiPrefix, routing.V1, routing.RouteRegister)
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
