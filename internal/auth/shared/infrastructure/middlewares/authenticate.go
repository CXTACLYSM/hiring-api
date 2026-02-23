package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/tokens"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/cache"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
	"github.com/golang-jwt/jwt/v5"
)

type Authenticate struct {
	userCache   *cache.UserCache
	findOneUser findOne.Handler
	secretKey   []byte
}

func NewAuthenticate(
	secretKey []byte,
	userCache *cache.UserCache,
	findOneUser findOne.Handler,
) *Authenticate {
	return &Authenticate{
		secretKey:   secretKey,
		userCache:   userCache,
		findOneUser: findOneUser,
	}
}

func (m *Authenticate) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := httputils.ExtractBearerToken(r)
		if !ok {
			httputils.ResponseError(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		userCacheKey := m.userCache.Key(token)
		if user, found := m.userCache.Get(userCacheKey); found {
			ctx := context.WithValue(r.Context(), middlewares.UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims := &tokens.Claims{}
		_, err := jwt.NewParser().ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.secretKey, nil
		})
		if err != nil {
			httputils.ResponseError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		user, err := m.findOneUser.Handle(findOne.Query{
			Id: claims.UserId,
		})
		if err != nil {
			log.Printf("error finding user by id=%s: %v", claims.UserId, err)
			httputils.ResponseError(w, http.StatusUnauthorized, "authentication failed")
			return
		}
		if user == nil {
			httputils.ResponseError(w, http.StatusUnauthorized, "user not found")
			return
		}

		err = m.userCache.Set(userCacheKey, user.ToShared())
		if err != nil {
			log.Printf("error SET redis key %s: %v", userCacheKey, err)
		}

		ctx := context.WithValue(r.Context(), middlewares.UserKey, user.ToShared())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
