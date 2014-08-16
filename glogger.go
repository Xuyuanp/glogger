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

// Logger is an interface supported method like Debug, Info and so on
type Logger interface {
	Leveler
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

var loggerRegister = NewRegister()

// GetLogger return a Logger with name.
// GetLogger will return nil if there is no Logger with this name
func GetLogger(name string) Logger {
	if v := loggerRegister.Get(name); v != nil {
		return v.(Logger)
	}
	return nil
}

// UnregisterLogger unregister the logger from global manager, this will make the logger
// unreachable for others. If this Logger hasn't been registered, nothing will happen.
func UnregisterLogger(name string) {
	loggerRegister.Unregister(name)
}

// RegisterLogger will register the logger to global manager. The logger registered can be
// accessed by GetLogger() method with logger's name.
func RegisterLogger(name string, l Logger) {
	loggerRegister.Register(name, l)
}
