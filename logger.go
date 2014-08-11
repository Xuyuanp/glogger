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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"
)

// gLogger is the default Logger
type gLogger struct {
	GroupFilter
	handlerGroup
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
	RegisterLogger(l)
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
	if !l.Filter(rec) {
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

func (l *gLogger) SetName(name string) {
	UnRegisterLogger(l)
	l.name = name
	RegisterLogger(l)
}

func (l *gLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *gLogger) LoadConfig(config []byte) {
	var m map[string]interface{}
	err := json.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}
	l.LoadConfigFromMap(m)
}

func (l *gLogger) LoadConfigFromMap(config map[string]interface{}) {
	name, ok := config["name"]
	if ok {
		l.SetName(name.(string))
	}
	level, ok := config["level"]
	if ok {
		l.level = StringToLevel[level.(string)]
	}
	handlers, ok := config["handlers"]
	if ok {
		for _, hname := range handlers.([]interface{}) {
			handler := GetHandler(hname.(string))
			l.AddHandler(handler)
		}
	}
}

func (l *gLogger) LoadConfigFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	code, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	l.LoadConfig(code)
}
