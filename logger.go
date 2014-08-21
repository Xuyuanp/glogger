/*
 * Copyright 2014 Xuyuan Pang <xuyuanp@gmail.com>
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

// gLogger is the default Logger
type gLogger struct {
	GroupFilter
	handlerGroup
	name  string
	level LogLevel
	ch    chan *Record
}

// NewLogger return a new Logger with debug level as default.
func NewLogger() *gLogger {
	l := &gLogger{
		level: DebugLevel,
		ch:    make(chan *Record, 100000),
	}
	return l
}

// Debug see details in Logger interface
func (l *gLogger) Debug(f string, v ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(f, v...))
}

// Info see details in Logger interface
func (l *gLogger) Info(f string, v ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(f, v...))
}

// Warning see details in Logger interface
func (l *gLogger) Warning(f string, v ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(f, v...))
}

// Error see details in Logger interface
func (l *gLogger) Error(f string, v ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(f, v...))
}

// Critical see details in Logger interface
func (l *gLogger) Critical(f string, v ...interface{}) {
	l.log(CriticalLevel, fmt.Sprintf(f, v...))
}

func (l *gLogger) log(level LogLevel, msg string) {
	if level < l.level {
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
	rec := NewRecord(l.name, now, level, file, funcname, line, msg)
	if !l.Filter(rec) {
		return
	}
	l.Handle(rec)
}

func (l *gLogger) run() {
	for {
		select {
		case rec := <-l.ch:
			l.Handle(rec)
		}
	}
}

func (l *gLogger) Level() LogLevel {
	return l.level
}

func (l *gLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *gLogger) LoadConfig(config map[string]interface{}) {
	if level, ok := config["level"]; ok {
		l.level = StringToLevel[level.(string)]
	} else {
		panic("'level' field is required")
	}
	if handlers, ok := config["handlers"]; ok {
		if len(handlers.([]interface{})) == 0 {
			panic("handler name is required")
		}
		for _, hname := range handlers.([]interface{}) {
			if handler := GetHandler(hname.(string)); handler != nil {
				l.AddHandler(handler)
			} else {
				panic("unknown handler name: " + hname.(string))
			}
		}
	} else {
		panic("'handlers' field is required")
	}
}
