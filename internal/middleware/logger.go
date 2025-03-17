package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

// Logger es un middleware que registra información detallada sobre cada solicitud HTTP
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Registrar información básica de la solicitud
		log.Printf("Solicitud recibida: %s %s", r.Method, r.URL.Path)
		log.Printf("Headers: %v", r.Header)
		
		// Si es una solicitud POST o PUT, registrar el cuerpo
		if r.Method == "POST" || r.Method == "PUT" {
			// Leer el cuerpo
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = io.ReadAll(r.Body)
				// Restaurar el cuerpo para que pueda ser leído nuevamente
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
			
			// Registrar el cuerpo de la solicitud
			log.Printf("Cuerpo de la solicitud: %s", string(bodyBytes))
		}
		
		// Procesar la solicitud
		next.ServeHTTP(w, r)
		
		// Registrar información sobre la respuesta
		log.Printf("Solicitud completada: %s %s en %v", r.Method, r.URL.Path, time.Since(start))
	})
}