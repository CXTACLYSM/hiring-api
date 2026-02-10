package middlewares

import "net/http"

type ContentType struct {
	contentType string
}

func NewContentType(contentType string) *ContentType {
	return &ContentType{
		contentType: contentType,
	}
}

func (m *ContentType) ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
