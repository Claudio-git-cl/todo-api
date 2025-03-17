package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/gorilla/mux"
	"github.com/claudio/todo-api/internal/models"
)

func TestCreateTask(t *testing.T) {
	// Configurar la prueba
	task := models.Task{
		Title:       "Tarea de prueba",
		Description: "Descripción de prueba",
		Completed:   false,
	}
	
	// Convertir la tarea a JSON
	taskJSON, _ := json.Marshal(task)
	
	// Crear una solicitud HTTP de prueba
	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}
	
	// Establecer el encabezado Content-Type
	req.Header.Set("Content-Type", "application/json")
	
	// Crear un ResponseRecorder para registrar la respuesta
	rr := httptest.NewRecorder()
	
	// Crear un router con el handler
	router := mux.NewRouter()
	taskHandler := NewTaskHandler()
	router.HandleFunc("/api/tasks", taskHandler.CreateTask).Methods("POST")
	
	// Ejecutar la solicitud
	router.ServeHTTP(rr, req)
	
	// Verificar el código de estado
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("El handler devolvió un código de estado incorrecto: obtuvo %v, esperaba %v",
			status, http.StatusCreated)
	}
	
	// Verificar el cuerpo de la respuesta
	var responseTask models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &responseTask)
	if err != nil {
		t.Fatal(err)
	}
	
	// Verificar que la tarea tenga un ID asignado
	if responseTask.ID == 0 {
		t.Error("Se esperaba que la tarea tuviera un ID asignado")
	}
	
	// Verificar que los campos coincidan
	if responseTask.Title != task.Title {
		t.Errorf("Título incorrecto: obtuvo %v, esperaba %v", responseTask.Title, task.Title)
	}
	
	if responseTask.Description != task.Description {
		t.Errorf("Descripción incorrecta: obtuvo %v, esperaba %v", responseTask.Description, task.Description)
	}
	
	if responseTask.Completed != task.Completed {
		t.Errorf("Estado completado incorrecto: obtuvo %v, esperaba %v", responseTask.Completed, task.Completed)
	}
}

func TestGetTasks(t *testing.T) {
	// Crear una solicitud HTTP de prueba
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Crear un ResponseRecorder para registrar la respuesta
	rr := httptest.NewRecorder()
	
	// Crear un router con el handler
	router := mux.NewRouter()
	taskHandler := NewTaskHandler()
	router.HandleFunc("/api/tasks", taskHandler.GetTasks).Methods("GET")
	
	// Ejecutar la solicitud
	router.ServeHTTP(rr, req)
	
	// Verificar el código de estado
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("El handler devolvió un código de estado incorrecto: obtuvo %v, esperaba %v",
			status, http.StatusOK)
	}
	
	// Verificar que la respuesta sea un array JSON
	var tasks []models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatal(err)
	}
}