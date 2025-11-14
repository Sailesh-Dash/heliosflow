package middleware

import (
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/Sailesh-Dash/heliosflow/internal/logger"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Recoverer just forwards to chi's Recoverer.
func Recoverer(next http.Handler) http.Handler {
	return chimw.Recoverer(next)
}
