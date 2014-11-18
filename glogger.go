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

// LogLevel type
type LogLevel uint8

// AutoRoot is a switcher that controlls GetLogger result
var AutoRoot = true

// LogLevel values
const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
)

// LevelToString is a map to translate LogLevel to a level name string
var LevelToString = map[LogLevel]string{
	DebugLevel:    "DBUG",
	InfoLevel:     "INFO",
	WarnLevel:     "WARN",
	ErrorLevel:    "ERRO",
	CriticalLevel: "CRIT",
}

// StringToLevel is a map to translate level name to LogLevel type
var StringToLevel = map[string]LogLevel{
	"DEBUG":    DebugLevel,
	"INFO":     InfoLevel,
	"WARNING":  WarnLevel,
	"ERROR":    ErrorLevel,
	"CRITICAL": CriticalLevel,
}

// Leveler is an interface provided set/get LogLevel method
type Leveler interface {
	Level() LogLevel
	SetLevel(level LogLevel)
}

var loggerRegister = NewRegister()

// GetLogger return a Logger registered with the name,
// or the root Logger if AutoRoot is true and name is not root, or nil
func GetLogger(name string) *Logger {
	if v := loggerRegister.Get(name); v != nil {
		return v.(*Logger)
	}
	if AutoRoot && name != "root" {
		return loggerRegister.Get("root").(*Logger)
	}
	return nil
}

// UnregisterLogger unregister the logger from global manager, this will make the logger
// unreachable for others and return the logger, this's the last chance getting it.
// If this Logger hasn't been registered, nothing will happen and return nil
func UnregisterLogger(name string) *Logger {
	if l := loggerRegister.Unregister(name); l != nil {
		return l.(*Logger)
	}
	return nil
}

// RegisterLogger will register the logger to global manager. The logger registered can be
// accessed by GetLogger() method with logger's name.
func RegisterLogger(name string, l *Logger) {
	loggerRegister.Register(name, l)
	l.Name = name
}
