// Package logger provides logger for the project
package logger

import (
	"log"
	"os"
)

// Create a new logger with the prefix given
func NewLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix+": ", log.Ltime)
}
