package logger

import (
	"log"
	"os"
)

// Global loggers available to the entire application.
var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

// init sets up the loggers once when the package is imported.
func init() {
	Info = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	Warn = log.New(os.Stderr, "WARN\t", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
}
