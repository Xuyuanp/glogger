package glogger

import "fmt"

var LevelMap = map[LogLevel]string{
	DebugLevel:    "DEBUG",
	InfoLevel:     "INFO",
	WarnLevel:     "WARNING",
	CriticalLevel: "CRITICAL",
}

type Formatter interface {
	Format(rec *Record) string
}

type DefaultFormatter struct {
	Fmt string
}

func NewDefaultFormatter(format string) Formatter {
	if format == "" {
		format = "[%s\t%s\t%s\t:%d] %s"
	}
	df := &DefaultFormatter{
		Fmt: format,
	}
	return df
}

func (df *DefaultFormatter) Format(rec *Record) string {
	levelName, _ := LevelMap[rec.Level]
	return fmt.Sprintf(df.Fmt, rec.Time, levelName, rec.File, rec.Line, rec.Message)
}
