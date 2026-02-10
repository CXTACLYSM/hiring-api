package httputils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/shared/application/response"
	"github.com/CXTACLYSM/hiring-api/internal/auth/tokens"
	appErrors "github.com/CXTACLYSM/hiring-api/internal/auth/user/application/errors"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/middlewares"
)

func WriteError(w http.ResponseWriter, err error) {
	var appError *appErrors.AppError
	if errors.As(err, &appError) {
		ResponseError(w, appError.StatusCode, appError.Error())
	} else {
		log.Printf("internal server error: %v", err)
		ResponseError(w, http.StatusInternalServerError, appErrors.MessageInternalServerError)
	}
}

func DecodeJSON(r *http.Request, dest any) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

func ResponseSuccess(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.SuccessResponse{
		Ok:   true,
		Data: data,
	})
}

func ResponseError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.ErrorResponse{
		Ok:      false,
		Message: message,
	})
}

func ClaimsFromRequest(r *http.Request) *tokens.Claims {
	return r.Context().Value(middlewares.ClaimsKey).(*tokens.Claims)
}
