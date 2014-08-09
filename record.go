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

type Record struct {
	Name    string
	Level   LogLevel
	Time    time.Time
	LFile   string
	SFile   string
	Line    int
	Func    string
	Message string
}

var pathReg = regexp.MustCompile("/.*/")

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
