package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/claudio/todo-api/internal/logger"
	"github.com/claudio/todo-api/internal/router"
	"github.com/gorilla/mux"
)

func main() {
	// Inicializar el logger
	logger.Init()
	logger.InfoLogger.Println("Iniciando la aplicación...")

	// Esperar un momento
	time.Sleep(1 * time.Second)

	logger.InfoLogger.Println("Utilizando almacenamiento en memoria para desarrollo...")
	
	// Inicializar el router sin base de datos (usará almacenamiento en memoria)
	r := router.NewRouter(nil)
	logger.InfoLogger.Println("Router inicializado correctamente")
	
	// Agregar registro de rutas para depuración
	routesList := []string{}
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		routesList = append(routesList, fmt.Sprintf("Ruta: %s, Métodos: %s", path, strings.Join(methods, ", ")))
		return nil
	})
	logger.InfoLogger.Printf("Rutas registradas:\n%s", strings.Join(routesList, "\n"))

	// Obtener el puerto del entorno o usar el predeterminado
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Implementar un middleware CORS más robusto
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Registrar detalles de la solicitud para depuración
			logger.InfoLogger.Printf("Solicitud recibida: %s %s", r.Method, r.URL.Path)
			logger.InfoLogger.Printf("Headers: %v", r.Header)
			
			// Configurar encabezados CORS para todas las respuestas
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Max-Age", "3600")
			
			// Manejar solicitudes preflight OPTIONS
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				logger.InfoLogger.Printf("Respondiendo a solicitud OPTIONS con 200 OK")
				return
			}
			
			// Procesar la solicitud normal
			h.ServeHTTP(w, r)
		})
	}

	// Configurar el servidor HTTP con el middleware CORS mejorado
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsHandler(r),
		// Aumentar los tiempos de espera para evitar problemas de conexión
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar el servidor
	logger.InfoLogger.Printf("Servidor iniciado en el puerto %s", port)
	if err := server.ListenAndServe(); err != nil {
		logger.ErrorLogger.Fatalf("Error al iniciar el servidor: %v", err)
	}
}