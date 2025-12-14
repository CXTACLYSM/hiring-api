package middleware

import (
	"net/http"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if true {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
		//ctx := context.WithValue(r.Context(), userKey, user)
		//next.ServeHTTP(w, r.WithContext(ctx))
	})
}
