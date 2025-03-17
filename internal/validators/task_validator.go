package validators

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/claudio/todo-api/internal/models"
)

// ValidateTask valida los datos de una tarea
func ValidateTask(r *http.Request) (*models.Task, error) {
	var task models.Task
	
	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return nil, err
	}
	
	// Validar que el título no esté vacío
	if task.Title == "" {
		return nil, errors.New("el título no puede estar vacío")
	}
	
	// Validar longitud máxima del título
	if len(task.Title) > 100 {
		return nil, errors.New("el título no puede tener más de 100 caracteres")
	}
	
	return &task, nil
}