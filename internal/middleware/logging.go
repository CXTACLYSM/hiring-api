package middleware

import (
	"log"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if true {
			log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		}

		next.ServeHTTP(w, r)
		//ctx := context.WithValue(r.Context(), userKey, user)
		//next.ServeHTTP(w, r.WithContext(ctx))
	})
}
