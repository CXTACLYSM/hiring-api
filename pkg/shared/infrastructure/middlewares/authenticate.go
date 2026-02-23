package middlewares

import (
	"context"
	"log"
	"net/http"

	pb "github.com/CXTACLYSM/hiring-api/pkg/grpc/auth/v1"
	entities "github.com/CXTACLYSM/hiring-api/pkg/shared/domain"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/cache"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
)

type contextKey string

const UserKey contextKey = "user"

type Authenticate struct {
	userCache *cache.UserCache
	client    pb.AuthServiceClient
}

func NewAuthenticate(userCache *cache.UserCache, client pb.AuthServiceClient) *Authenticate {
	return &Authenticate{
		userCache: userCache,
		client:    client,
	}
}

func (m *Authenticate) Authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token, ok := httputils.ExtractBearerToken(r)
		if !ok {
			httputils.ResponseError(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		userCacheKey := m.userCache.Key(token)
		user, ok := m.userCache.Get(userCacheKey)
		if ok && user != nil {
			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		resp, err := m.client.ValidateToken(r.Context(), &pb.ValidateTokenRequest{
			Token: token,
		})
		if err != nil {
			httputils.ResponseError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		if resp.User == nil {
			httputils.ResponseError(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		user = &entities.User{
			Id:       resp.User.Id,
			Username: resp.User.Username,
			Email:    resp.User.Email,
		}

		err = m.userCache.Set(userCacheKey, user)
		if err != nil {
			log.Printf("error SET redis key %s: %v", userCacheKey, err)
		}

		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func UserFromRequest(r *http.Request) *entities.User {
	user, ok := r.Context().Value(UserKey).(*entities.User)
	if !ok {
		return nil
	}
	return user
}
