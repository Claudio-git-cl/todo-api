package logger

import (
	"log"
	"os"
)

var (
	// InfoLogger registra información general
	InfoLogger *log.Logger
	// ErrorLogger registra errores
	ErrorLogger *log.Logger
)

// Init inicializa los loggers
func Init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}