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
	logFlags := log.Lshortfile | log.Ltime | log.Lmsgprefix

	debug := log.New(os.Stdout, "DEBUG ", logFlags)
	info := log.New(os.Stdout, "INFO ", logFlags)
	warning := log.New(os.Stdout, "WARNING ", logFlags)
	err := log.New(os.Stderr, "ERROR ", logFlags)

	return &Logger{
		level:         level,
		debugLogger:   debug,
		infoLogger:    info,
		warningLogger: warning,
		errorLogger:   err,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.infoLogger.Printf(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	if l.level > Debug {
		return
	}
	l.debugLogger.Printf(msg, args...)
}
