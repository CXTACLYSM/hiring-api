package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CXTACLYSM/hiring-api/configs"
)

type InfoHandler struct {
	cfg *configs.Config
}

func NewInfoHandler(cfg *configs.Config) *InfoHandler {
	return &InfoHandler{
		cfg: cfg,
	}
}

func (h *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(map[string]interface{}{
		"version": h.cfg.App.Version,
	})
	_, err = w.Write(body)
	if err != nil {
		fmt.Printf("error writing response: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}
