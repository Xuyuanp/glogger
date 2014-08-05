package glogger

import "time"

type Record struct {
	Level   LogLevel
	Time    time.Time
	File    string
	Line    int
	Message string
}

func NewRecord(t time.Time, level LogLevel, file string, line int, msg string) *Record {
	rec := &Record{
		Level:   level,
		Time:    t,
		File:    file,
		Line:    line,
		Message: msg,
	}
	return rec
}
