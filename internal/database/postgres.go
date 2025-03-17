package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/claudio/todo-api/internal/logger"
	_ "github.com/lib/pq"
)

// NewPostgresConnection establece una conexión con la base de datos PostgreSQL
func NewPostgresConnection() (*sql.DB, error) {
	// Obtener variables de entorno para la conexión
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "todo_db")

	logger.InfoLogger.Printf("Conectando a la base de datos: %s:%s/%s", host, port, dbname)

	// Construir la cadena de conexión
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Abrir la conexión
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.ErrorLogger.Printf("Error al abrir la conexión a la base de datos: %v", err)
		return nil, err
	}

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		logger.ErrorLogger.Printf("Error al verificar la conexión a la base de datos: %v", err)
		return nil, err
	}

	logger.InfoLogger.Println("Conexión a la base de datos establecida correctamente")
	return db, nil
}

// getEnv obtiene una variable de entorno o devuelve un valor predeterminado
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}