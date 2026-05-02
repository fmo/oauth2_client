// Package internal
package internal

import (
	"log"
	"os"
)

type Level int

const (
	Debug = iota
	Info
	Warning
	Error
)

type Logger struct {
	level Level

	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func NewLogger(level Level) *Logger {
	debug := log.New(os.Stdout, "DEBUG ", log.Llongfile)
	info := log.New(os.Stdout, "INFO ", log.Llongfile)
	warning := log.New(os.Stdout, "WARNING ", log.Llongfile)
	err := log.New(os.Stderr, "ERROR ", log.Llongfile)

	return &Logger{
		level:         level,
		debugLogger:   debug,
		infoLogger:    info,
		warningLogger: warning,
		errorLogger:   err,
	}
}

func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
}
