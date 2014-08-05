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
	Filterer
	Handlers []Handler
	Name     string
	Level    LogLevel
	Parent   *Logger
}

func NewLogger(name string, level LogLevel) *Logger {
	l := &Logger{
		Name:     name,
		Level:    level,
		Handlers: []Handler{},
	}
	return l
}

func (l *Logger) AddHandler(h Handler) {
	if len(l.Handlers) == 0 {
		l.Handlers = []Handler{h}
	} else {
		l.Handlers = append(l.Handlers[:len(l.Handlers)], h)
	}
}

func (l *Logger) Info(v ...interface{}) {
	var msg string
	if len(v) > 1 {
		msg = fmt.Sprintf(v[0].(string), v[1:])
	} else {
		msg = v[0].(string)
	}
	l.log(InfoLevel, msg)
}

func (l *Logger) log(level LogLevel, msg string) {
	if level < l.Level {
		return
	}
	now := time.Now()
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	rec := NewRecord(now, level, file, line, msg)
	if !l.Filter(rec) {
		return
	}
	for _, hdl := range l.Handlers {
		hdl.Handle(rec)
	}
}
