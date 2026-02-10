package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CXTACLYSM/hiring-api/internal/auth/shared/application/response"
	"github.com/CXTACLYSM/hiring-api/internal/auth/tokens"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClaimsKey contextKey = "claims"

type Authenticate struct {
	secretKey []byte
}

func NewAuthenticate(secretKey []byte) *Authenticate {
	return &Authenticate{secretKey: secretKey}
}

func (m *Authenticate) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.ErrorResponse{
				Ok:      false,
				Message: "missing or invalid authorization header",
			})
			return
		}
		encodedToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &tokens.Claims{}
		_, err := jwt.NewParser().ParseWithClaims(encodedToken, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.secretKey, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.ErrorResponse{
				Ok:      false,
				Message: "invalid or expired token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
