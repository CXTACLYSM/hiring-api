package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RegisterHandler struct {
	writePool *pgxpool.Pool
}

func NewRegisterHandler(writePool *pgxpool.Pool) *RegisterHandler {
	return &RegisterHandler{
		writePool: writePool,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
