/*
 * Copyright 2014 Xuyuan Pang <xuyuanp # gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package glogger

import (
	"fmt"
	"runtime"
	"time"
)

// Logger is the default Logger
type Logger struct {
	GroupFilter
	handlerGroup
	Name   string
	Level  LogLevel
	ch     chan *Record
	parent *Logger
}

// NewLogger return a new Logger with debug level as default.
func NewLogger() *Logger {
	l := &Logger{
		Level: DebugLevel,
	}
	return l
}

// Default functions return a logger registered by name 'root',
// or a new Logger with default Handler and Formatter, and registered as 'root' automatically.
func Default() *Logger {
	if l := GetLogger("root"); l != nil {
		return l
	}
	l := NewLogger()
	h := NewStreamHandler()
	l.AddHandler(h)
	RegisterLogger("root", l)
	return l
}

var std = Default()

// Debug see details in Logger interface
func (l *Logger) Debug(f string, v ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(f, v...))
}

// Info see details in Logger interface
func (l *Logger) Info(f string, v ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(f, v...))
}

// Warning see details in Logger interface
func (l *Logger) Warning(f string, v ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(f, v...))
}

// Error see details in Logger interface
func (l *Logger) Error(f string, v ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(f, v...))
}

// Critical see details in Logger interface
func (l *Logger) Critical(f string, v ...interface{}) {
	l.log(CriticalLevel, fmt.Sprintf(f, v...))
}

func (l *Logger) log(level LogLevel, msg string) {
	if level < l.Level {
		return
	}
	now := time.Now()
	pc, file, line, ok := runtime.Caller(2)
	var funcname string
	if !ok {
		file = "???"
		line = 0
		funcname = "???"
	} else {
		funcname = runtime.FuncForPC(pc).Name()
	}
	rec := NewRecord(l.Name, now, level, file, funcname, line, msg)
	if !l.Filter(rec) {
		return
	}
	l.Handle(rec)
}

func (l *Logger) run() {
	for {
		select {
		case rec := <-l.ch:
			l.Handle(rec)
		}
	}
}

// LoadConfig loads configuration from map.
func (l *Logger) LoadConfig(config map[string]interface{}) error {
	// Load log level, default is DebugLevel
	if blevel, ok := config["level"]; ok {
		if level, ok := StringToLevel[blevel.(string)]; ok {
			l.Level = level
		} else {
			return fmt.Errorf("unknown log level: %s", blevel.(string))
		}
	} else {
		l.Level = DebugLevel
	}
	// Load handlers, default is StreamHandler
	if handlers, ok := config["handlers"]; ok && len(handlers.([]interface{})) > 0 {
		l.ClearHandlers()
		for _, hname := range handlers.([]interface{}) {
			if handler := GetHandler(hname.(string)); handler != nil {
				l.AddHandler(handler)
			} else {
				return fmt.Errorf("unknown handler name: %s", hname.(string))
			}
		}
	} else {
		h := NewStreamHandler()
		l.SetHandlers(h)
	}
	// Load filters
	if filters, ok := config["filters"]; ok {
		if len(filters.([]interface{})) > 0 {
			for _, fname := range filters.([]interface{}) {
				if filter := GetFilter(fname.(string)); filter != nil {
					l.AddFilter(filter)
				} else {
					return fmt.Errorf("unknown filter name: %s", fname.(string))
				}
			}
		}
	}
	return nil
}
