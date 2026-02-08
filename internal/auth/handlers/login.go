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

type LoginHandler struct {
	middlewares []middlewares.Middleware
	authService *services.AuthService
}

func NewLoginHandler(authService *services.AuthService, middlewares ...middlewares.Middleware) *LoginHandler {
	return &LoginHandler{
		authService: authService,
		middlewares: middlewares,
	}
}

func (h *LoginHandler) Middlewares() []middlewares.Middleware {
	return h.middlewares
}

func (h *LoginHandler) RouteMetadata() *routing.Metadata {
	return routing.NewRouteMetadata(routing.ApiPrefix, routing.V1, routing.RouteLogin)
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Ok:      false,
			Message: "method not allowed",
		})
		return
	}

	var loginDTO *dto.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Ok:      false,
			Message: "invalid request body",
		})
		return
	}

	var appError *appErrors.AppError
	token, err := h.authService.Login(loginDTO)
	if err != nil {
		if errors.As(err, &appError) {
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Ok:      false,
				Message: appError.Error(),
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.SuccessResponse{
		Ok:   true,
		Data: map[string]string{"token": token},
	})
}
