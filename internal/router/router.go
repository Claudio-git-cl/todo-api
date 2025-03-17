package router

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/claudio/todo-api/internal/handlers"
	"github.com/claudio/todo-api/internal/middleware"
	"github.com/claudio/todo-api/internal/logger"
)

// NewRouter configura y devuelve un nuevo router
func NewRouter(db *sql.DB) *mux.Router {
	logger.InfoLogger.Println("Configurando router...")
	
	r := mux.NewRouter()

	// Aplicar middleware CORS a todas las rutas
	r.Use(middleware.CORSMiddleware)
	// Aplicar middleware de registro a todas las rutas
	r.Use(middleware.Logger)

	// Crear el manejador de tareas
	taskHandler := handlers.NewTaskHandler()

	// Endpoint de prueba
	r.HandleFunc("/api/health", taskHandler.HealthCheck).Methods("GET")

	// Definir las rutas
	r.HandleFunc("/api/tasks", taskHandler.GetTasks).Methods("GET")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", taskHandler.GetTask).Methods("GET")
	
	// Ruta para crear tareas - asegurarse de que esté correctamente configurada
	logger.InfoLogger.Println("Configurando ruta POST para crear tareas")
	r.HandleFunc("/api/tasks", taskHandler.CreateTask).Methods("POST")
	
	r.HandleFunc("/api/tasks/{id:[0-9]+}", taskHandler.UpdateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", taskHandler.DeleteTask).Methods("DELETE")
	
	// Agregar manejo de solicitudes OPTIONS para CORS
	logger.InfoLogger.Println("Configurando rutas OPTIONS para CORS")
	r.HandleFunc("/api/tasks", taskHandler.HandlePreflight).Methods("OPTIONS")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", taskHandler.HandlePreflight).Methods("OPTIONS")

	// Configurar ruta para manejar todas las solicitudes OPTIONS (para mayor seguridad)
	r.PathPrefix("/").HandlerFunc(taskHandler.HandlePreflight).Methods("OPTIONS")

	logger.InfoLogger.Println("Router configurado correctamente")
	return r
}