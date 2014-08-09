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
	"sync"
	"time"
)

// LogLevel type
type LogLevel uint8

// log message level
const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
)

// Logger is an interface supported functions like Debug, Info and so on
type Logger interface {
	Filter

	// log DebugLevel message
	Debug(f string, v ...interface{})

	// log InfoLevel message
	Info(f string, v ...interface{})

	// log WarnLevel message
	Warning(f string, v ...interface{})

	// log ErrorLevel message
	Error(f string, v ...interface{})

	// log CriticalLevel message
	Critical(f string, v ...interface{})

	// Name return the name of Logger
	Name() string

	// Level return the LogLevel of Logger
	Level() LogLevel

	AddHandler(h Handler)
}

type loggerMapper struct {
	mapper map[string]Logger
	mu     sync.RWMutex
}

var lm *loggerMapper
var once sync.Once

func init() {
	// make sure loggerMapper init only once
	once.Do(setup)
}

func setup() {
	lm = &loggerMapper{
		mapper: map[string]Logger{},
	}
}

func (lm *loggerMapper) GetLogger(name string) Logger {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.mapper[name]
}

func (lm *loggerMapper) registerLogger(l Logger) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	_, ok := lm.mapper[l.Name()]
	if ok {
		panic(fmt.Sprintf("Register logger with name:%s twice", l.Name()))
	}
	lm.mapper[l.Name()] = l
}

// GetLogger return a Logger with name.
// GetLogger will return nil if there is no Logger with this name
func GetLogger(name string) Logger {
	return lm.GetLogger(name)
}

func registerLogger(l Logger) {
	lm.registerLogger(l)
}

// gLogger is the default Logger
type gLogger struct {
	GroupFilter
	handlerManger
	name  string
	level LogLevel
}

// New return a new Logger.
// name means the logger's name, it should be unique.
// level means the logger's level, all the logs who's level lower than this will be ignore
// It will panic if this name has been registered.
func New(name string, level LogLevel) Logger {
	l := &gLogger{
		name:  name,
		level: level,
	}
	registerLogger(l)
	return l
}

func (l *gLogger) Debug(f string, v ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(f, v...))
}

func (l *gLogger) Info(f string, v ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(f, v...))
}

func (l *gLogger) Warning(f string, v ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(f, v...))
}

func (l *gLogger) Error(f string, v ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(f, v...))
}

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
	if !l.DoFilter(rec) {
		return
	}
	l.Handle(rec)
}

func (l *gLogger) Name() string {
	return l.name
}

func (l *gLogger) Level() LogLevel {
	return l.level
}
