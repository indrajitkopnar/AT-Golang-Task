package utils

import (
	"log"
	"os"
)

// Logger instance
var Logger *log.Logger

// Initialize Logger
func InitLogger() {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	Logger = log.New(logFile, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
