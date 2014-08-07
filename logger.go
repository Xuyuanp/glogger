package glogger

import (
	"fmt"
	"runtime"
	"time"
)

type LogLevel int

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
)

type Logger struct {
	FilterGroup
	HandlerGroup
	Name   string
	Level  LogLevel
	Parent *Logger
}

func NewLogger(name string, level LogLevel) *Logger {
	l := &Logger{
		Name:  name,
		Level: level,
	}
	l.AddFilter(NewLevelFilter(level))
	return l
}

func (l *Logger) log(level LogLevel, msg string) {
	now := time.Now()
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	rec := NewRecord(l.Name, now, level, file, line, msg)
	if !l.DoFilter(rec) {
		return
	}
	l.Handle(rec)
}

func (l *Logger) Info(f string, v ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(f, v...))
}
