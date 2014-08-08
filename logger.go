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

type LogLevel int

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
)

type Logger struct {
	GroupFilter
	GroupHandler
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

func (l *Logger) Debug(f string, v ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(f, v...))
}

func (l *Logger) Warning(f string, v ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(f, v...))
}

func (l *Logger) Error(f string, v ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(f, v...))
}

func (l *Logger) Critical(f string, v ...interface{}) {
	l.log(CriticalLevel, fmt.Sprintf(f, v...))
}
