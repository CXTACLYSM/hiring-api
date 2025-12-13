package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	writePool *pgxpool.Pool
}

func NewAuthHandler(writePool *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		writePool: writePool,
	}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
