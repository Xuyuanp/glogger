package glogger

import "time"

type Record struct {
	Name    string
	Level   LogLevel
	Time    time.Time
	File    string
	Line    int
	Message string
}

func NewRecord(name string, t time.Time, level LogLevel, file string, line int, msg string) *Record {
	rec := &Record{
		Name:    name,
		Level:   level,
		Time:    t,
		File:    file,
		Line:    line,
		Message: msg,
	}
	return rec
}
