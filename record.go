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
	"regexp"
	"time"
)

// Record is a struct contains all the logging information
type Record struct {
	Name    string    // logger name
	Level   LogLevel  // log level
	Time    time.Time // log time
	LFile   string    // full file name
	SFile   string    // final file name
	Line    int       // line number
	Func    string    // function name
	Message string    // log message
}

var pathReg = regexp.MustCompile("/.*/")

// NewRecord return a new Record
func NewRecord(name string, t time.Time, level LogLevel, file string, funcname string, line int, msg string) *Record {
	rec := &Record{
		Name:    name,
		Level:   level,
		Time:    t,
		LFile:   file,
		Line:    line,
		Func:    funcname,
		Message: msg,
	}
	rec.SFile = pathReg.ReplaceAllString(file, "")
	return rec
}
