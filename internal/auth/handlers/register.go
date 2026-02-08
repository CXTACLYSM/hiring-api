package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/dto"
	appErrors "github.com/CXTACLYSM/hiring-api/internal/auth/errors"
	"github.com/CXTACLYSM/hiring-api/internal/auth/middlewares"
	"github.com/CXTACLYSM/hiring-api/internal/auth/routing"
	"github.com/CXTACLYSM/hiring-api/internal/auth/services"
)

type RegisterHandler struct {
	middlewares []middlewares.Middleware
	authService *services.AuthService
}

func NewRegisterHandler(authService *services.AuthService, middlewares ...middlewares.Middleware) *RegisterHandler {
	return &RegisterHandler{
		authService: authService,
		middlewares: middlewares,
	}
}

func (h *RegisterHandler) Middlewares() []middlewares.Middleware {
	return h.middlewares
}

func (h *RegisterHandler) RouteMetadata() *routing.Metadata {
	return routing.NewRouteMetadata(routing.ApiPrefix, routing.V1, routing.RouteRegister)
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Ok:      false,
			Message: "method not allowed",
		})
		return
	}

	var registerDTO *dto.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Ok:      false,
			Message: "invalid request body",
		})
		return
	}

	token, err := h.authService.Register(registerDTO)
	var appErr *appErrors.AppError
	if err != nil {
		if errors.As(err, &appErr) {
			w.WriteHeader(appErr.StatusCode)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Ok:      false,
				Message: appErr.Error(),
			})
		} else {
			log.Printf("internal server error: %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Ok:      false,
				Message: appErrors.MessageInternalServerError,
			})
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.SuccessResponse{
		Ok:   true,
		Data: map[string]string{"token": token},
	})
}
