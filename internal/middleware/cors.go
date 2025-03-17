package middleware

import (
	"net/http"
	"github.com/claudio/todo-api/internal/logger"
)

// CORSMiddleware maneja los encabezados CORS para todas las solicitudes
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log detailed request information
		logger.InfoLogger.Printf("CORS Middleware: Received %s request to %s", r.Method, r.URL.Path)
		
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			logger.InfoLogger.Printf("Handling OPTIONS preflight request for %s", r.URL.Path)
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Process the request
		next.ServeHTTP(w, r)
	})
}