package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware verifica que las solicitudes tengan un token válido
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token del encabezado Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Se requiere autorización", http.StatusUnauthorized)
			return
		}

		// Verificar el formato del token (Bearer token)
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Formato de autorización inválido", http.StatusUnauthorized)
			return
		}

		// Aquí implementarías la verificación real del token
		// Por ahora, simplemente permitimos cualquier token
		
		// Continuar con la siguiente función en la cadena
		next.ServeHTTP(w, r)
	})
}