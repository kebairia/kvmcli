package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string, keyAndValues ...any)
	Error(msg string, keyAndValues ...any)
}

// SimpleLogger is a basic implementation of the Logger interface.
type SimpleLogger struct {
	logger *log.Logger
}

// NewLogger creates a new SimpleLogger instance.
func NewLogger() Logger {
	return &SimpleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *SimpleLogger) Info(msg string, keysAndValues ...any) {
	l.logger.Printf("INFO: "+msg+" %v", keysAndValues)
}

func (l *SimpleLogger) Error(msg string, keysAndValues ...any) {
	l.logger.Printf("ERROR: "+msg+" %v", keysAndValues)
}
