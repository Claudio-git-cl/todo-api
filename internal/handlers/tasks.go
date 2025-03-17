package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/claudio/todo-api/internal/models"
)

// TaskHandler maneja las solicitudes relacionadas con tareas
type TaskHandler struct {
	tasks  []models.Task
	nextID int
}

// NewTaskHandler crea una nueva instancia de TaskHandler
func NewTaskHandler() *TaskHandler {
	handler := &TaskHandler{
		tasks:  []models.Task{},
		nextID: 1,
	}
	
	// Obtener la fecha y hora actual
	now := time.Now()
	
	// Agregar algunas tareas de ejemplo
	handler.tasks = append(handler.tasks, models.Task{
		ID:          handler.nextID,
		Title:       "Ejemplo de tarea 1",
		Description: "Esta es una tarea de ejemplo predefinida",
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	handler.nextID++
	
	handler.tasks = append(handler.tasks, models.Task{
		ID:          handler.nextID,
		Title:       "Ejemplo de tarea 2",
		Description: "Esta es otra tarea de ejemplo predefinida",
		Completed:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	handler.nextID++
	
	return handler
}

// HealthCheck proporciona un endpoint simple para verificar que la API está funcionando
func (h *TaskHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	response := map[string]string{
		"status": "ok",
		"message": "API funcionando correctamente",
	}
	
	json.NewEncoder(w).Encode(response)
}

// GetTasks devuelve todas las tareas
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	json.NewEncoder(w).Encode(h.tasks)
}

// GetTask devuelve una tarea específica por ID
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	
	// Obtener el ID de la URL
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	
	// Buscar la tarea
	for _, task := range h.tasks {
		if task.ID == id {
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	
	// Si no se encuentra la tarea
	http.Error(w, "Tarea no encontrada", http.StatusNotFound)
}

// CreateTask crea una nueva tarea
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	
	// Leer y registrar el cuerpo de la solicitud para depuración
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
		http.Error(w, "Error al leer la solicitud", http.StatusBadRequest)
		return
	}
	
	// Restaurar el cuerpo para que pueda ser leído por json.NewDecoder
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	
	// Registrar el cuerpo para depuración
	log.Printf("Cuerpo de la solicitud CreateTask: %s", string(bodyBytes))
	
	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("Error al decodificar JSON: %v", err)
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Validar los campos requeridos
	if task.Title == "" {
		log.Printf("Error: Título vacío")
		http.Error(w, "El título es obligatorio", http.StatusBadRequest)
		return
	}
	
	// Asignar un ID a la tarea y fechas
	task.ID = h.nextID
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	h.nextID++
	
	// Agregar la tarea a la lista
	h.tasks = append(h.tasks, task)
	
	// Registrar la tarea creada
	log.Printf("Tarea creada: %+v", task)
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// UpdateTask actualiza una tarea existente
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	
	// Obtener el ID de la URL
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	
	// Decodificar la tarea actualizada
	var updatedTask models.Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Buscar y actualizar la tarea
	for i, task := range h.tasks {
		if task.ID == id {
			// Mantener el ID original y la fecha de creación
			updatedTask.ID = id
			updatedTask.CreatedAt = task.CreatedAt
			updatedTask.UpdatedAt = time.Now()
			h.tasks[i] = updatedTask
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}
	
	// Si no se encuentra la tarea
	http.Error(w, "Tarea no encontrada", http.StatusNotFound)
}

// DeleteTask elimina una tarea
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	
	// Obtener el ID de la URL
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	
	// Buscar y eliminar la tarea
	for i, task := range h.tasks {
		if task.ID == id {
			// Eliminar la tarea del slice
			h.tasks = append(h.tasks[:i], h.tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	
	// Si no se encuentra la tarea
	http.Error(w, "Tarea no encontrada", http.StatusNotFound)
}

// HandlePreflight maneja las solicitudes OPTIONS para CORS
func (h *TaskHandler) HandlePreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.WriteHeader(http.StatusOK)
}