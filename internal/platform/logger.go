package platform

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetOutput(os.Stdout)

	return &Logger{logger}
}

func (l *Logger) Info(msg string) {
	l.Printf("[INFO] %s", msg)
}

func (l *Logger) Error(msg string) {
	l.Printf("[ERROR] %s", msg)
}
