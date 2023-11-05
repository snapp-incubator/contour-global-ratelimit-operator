package xdserver

import (
	"log"
)

func init() {
	logger = Logger{}
}

type Logger struct {
	Debug bool
}

func (logger Logger) Debugf(format string, args ...interface{}) {
	if logger.Debug {
		log.Printf("\t\tDEBUG\t"+format+"\n", args...)
	}
}

func (logger Logger) Infof(format string, args ...interface{}) {
	log.Printf("\t\tINFO\t"+format+"\n", args...)
}

func (logger Logger) Warnf(format string, args ...interface{}) {
	log.Printf("\t\tWARN\t"+format+"\n", args...)
}

func (logger Logger) Errorf(format string, args ...interface{}) {
	log.Printf("\t\tERROR\t"+format+"\n", args...)
}
