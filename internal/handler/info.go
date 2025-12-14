package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CXTACLYSM/hiring-api/configs"
	"github.com/CXTACLYSM/hiring-api/internal/middleware"
	"github.com/CXTACLYSM/hiring-api/internal/routing"
)

type InfoHandler struct {
	middlewares []middleware.Middleware
	cfg         *configs.Config
}

func NewInfoHandler(cfg *configs.Config, middlewares []middleware.Middleware) *InfoHandler {
	return &InfoHandler{
		cfg:         cfg,
		middlewares: middlewares,
	}
}

func (h *InfoHandler) RouteMetadata() *routing.Metadata {
	return routing.NewRouteMetadata(routing.ApiPrefix, routing.V1, routing.RouteInfo)
}

func (h *InfoHandler) Middlewares() []middleware.Middleware {
	return h.middlewares
}

func (h *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, err := json.Marshal(map[string]interface{}{
		"version": h.cfg.App.Version,
	})
	_, err = w.Write(body)
	if err != nil {
		fmt.Printf("error writing response: %v", err)
	}
}
