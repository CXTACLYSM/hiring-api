package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	pkgErrors "github.com/CXTACLYSM/hiring-api/pkg/shared/application/errors"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"github.com/go-playground/validator/v10"
)

func WriteError(w http.ResponseWriter, err error) {
	var appError *pkgErrors.ApplicationError
	var validationErrors validator.ValidationErrors

	if errors.As(err, &appError) {
		ResponseError(w, http.StatusUnprocessableEntity, appError.Error())
	} else if errors.As(err, &validationErrors) {
		ResponseValidationErrors(w, validationErrors)
	} else {
		log.Printf("internal server error: %v", err)
		ResponseError(w, http.StatusInternalServerError, "internal server error")
	}
}

func DecodeJSON(r *http.Request, dest any) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dest); err != nil {
		if errors.Is(err, io.EOF) {
			return &pkgErrors.ApplicationError{
				Message: "request body is required",
			}
		}

		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return &pkgErrors.ApplicationError{
				Message: "invalid json syntax",
			}
		}

		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			return &pkgErrors.ApplicationError{
				Message: fmt.Sprintf("invalid type for field %s: expected %s", typeErr.Field, typeErr.Type.String()),
			}
		}

		return &pkgErrors.ApplicationError{
			Message: "invalid request body",
		}
	}

	return nil
}

func ResponseOk(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ResponseOkNoData(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func ResponseError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(responses.ErrorResponse{
		Message: message,
	})
}

func ResponseValidationErrors(w http.ResponseWriter, errors validator.ValidationErrors) {
	w.WriteHeader(http.StatusUnprocessableEntity)

	json.NewEncoder(w).Encode(responses.ValidationErrorsResponse{
		Message: "validation errors",
		Errors:  validation.FormatErrors(errors),
	})
}

func ExtractBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", false
	}
	return strings.TrimPrefix(authHeader, "Bearer "), true
}
