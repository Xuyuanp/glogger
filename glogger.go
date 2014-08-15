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
	"sync"
)

// LogLevel type
type LogLevel uint8

// LogLevel values
const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
)

// LevelToString is a map to translat LogLevel to a level name string
var LevelToString = map[LogLevel]string{
	DebugLevel:    "DBUG",
	InfoLevel:     "INFO",
	WarnLevel:     "WARN",
	ErrorLevel:    "ERRO",
	CriticalLevel: "CRIT",
}

// StringToLevel is a map to translat level name to LogLevel type
var StringToLevel = map[string]LogLevel{
	"DEBUG":    DebugLevel,
	"INFO":     InfoLevel,
	"WARNING":  WarnLevel,
	"ERROR":    ErrorLevel,
	"CRITICAL": CriticalLevel,
}

// Namer is an interface provided set/get name method
type Namer interface {
	Name() string
	SetName(name string)
}

// Leveler is an interface provided set/get LogLevel method
type Leveler interface {
	Level() LogLevel
	SetLevel(level LogLevel)
}

// Logger is an interface supported method like Debug, Info and so on
type Logger interface {
	Leveler
	Namer
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

	AddHandler(h Handler)
}

type loggerManager struct {
	mapper map[string]Logger
	mu     sync.RWMutex
}

var lm *loggerManager
var once sync.Once

func init() {
	// make sure loggerManager init only once
	once.Do(setup)
}

func setup() {
	lm = &loggerManager{
		mapper: make(map[string]Logger),
	}
}

func (lm *loggerManager) getLogger(name string) Logger {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.mapper[name]
}

func (lm *loggerManager) unregisterLogger(l Logger) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	delete(lm.mapper, l.Name())
}

func (lm *loggerManager) unregisterLoggerByName(name string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	delete(lm.mapper, name)
}

func (lm *loggerManager) registerLogger(l Logger) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if _, ok := lm.mapper[l.Name()]; ok {
		panic(fmt.Sprintf("Register logger with name:%s twice", l.Name()))
	}
	lm.mapper[l.Name()] = l
}

// GetLogger return a Logger with name.
// GetLogger will return nil if there is no Logger with this name
func GetLogger(name string) Logger {
	return lm.getLogger(name)
}

// UnregisterLogger unregister the logger from global manager, this will make the logger
// unreachable for others. If this Logger hasn't been registered, nothing will happen.
func UnregisterLogger(l Logger) {
	lm.unRegisterLogger(l)
}

// UnregisterLoggerByName is like UnregisterLogger, but unregister the Logger by name.
// If there is no Logger with this name registered before, nothing will happen.
func UnregisterLoggerByName(name string) {
	lm.unRegisterLoggerByName(name)
}

// RegisterLogger will register the logger to global manager. The logger registered can be
// accessed by GetLogger() method with logger's name.
func RegisterLogger(l Logger) {
	lm.registerLogger(l)
}
