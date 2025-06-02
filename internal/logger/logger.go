// logger/logger.go
package logger

import (
	"fmt"
	"log/slog"
	"os"
)

var std *slog.Logger

func init() {
	// Create a TextHandler that writes to stdout.
	// We’re not using any special HandlerOptions here,
	// so Debug, Info, etc. all come through as simple text.
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// By default, TextHandler prints with a timestamp and level.
		// If you ever want to add file:line, you can set AddSource: true.
	})
	std = slog.New(handler)
}

// Debug logs a plain message at LevelDebug (no printf formatting).
func Debug(msg string) {
	std.Debug(msg)
}

// Debugf formats the string and logs at LevelDebug.
func Debugf(format string, args ...any) {
	std.Debug(fmt.Sprintf("\033[35mdebug: \033[0m"+format, args...))
}

// Info logs a plain message at LevelInfo.
func Info(msg string) {
	std.Info(msg)
}

// Infof formats the string and logs at LevelInfo.
func Infof(format string, args ...any) {
	std.Info(fmt.Sprintf(format, args...))
}

// Error prints a plain, unstructured error message to stderr.
// It does NOT go through slog—no timestamps, no level prefixes.
func Error(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

// Errorf formats and prints an unstructured error message to stderr.
func Errorf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "\033[31merror: \033[0m"+format+"\n", args...)
}

// Warn prints a warning (unstructured) to stderr, in yellow.
func Warn(msg string) {
	fmt.Fprintln(os.Stderr, "\033[33mwarning:\033[0m", msg)
}

func Warnf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "\033[33mwarning:\033[0m "+format+"\n", args...)
}

// Fatal prints an error then exits with code 1.
func Fatal(msg string) {
	fmt.Fprintln(os.Stderr, "\033[31mfatal:\033[0m", msg)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "\033[31mfatal:\033[0m "+format+"\n", args...)
	os.Exit(1)
}
