package xdserver

import (
	"log"
)

func init() {
	logger = Logger{}
}

// Logger is a custom logger with different log levels.
type Logger struct {
	Debug bool
}

// Debugf logs a debug message if Debug is true.
func (logger Logger) Debugf(format string, args ...interface{}) {
	if logger.Debug {
		log.Printf("\t\tDEBUG\t"+format+"\n", args...)
	}
}

// Infof logs an informational message.
func (logger Logger) Infof(format string, args ...interface{}) {
	log.Printf("\t\tINFO\t"+format+"\n", args...)
}

// Warnf logs a warning message.
func (logger Logger) Warnf(format string, args ...interface{}) {
	log.Printf("\t\tWARN\t"+format+"\n", args...)
}

// Errorf logs an error message.
func (logger Logger) Errorf(format string, args ...interface{}) {
	log.Printf("\t\tERROR\t"+format+"\n", args...)
}
