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
	"regexp"
	"strings"
	"time"
)

var LevelMap = map[LogLevel]string{
	DebugLevel:    "DBUG",
	InfoLevel:     "INFO",
	WarnLevel:     "WARN",
	ErrorLevel:    "ERRO",
	CriticalLevel: "CRIT",
}

type Formatter interface {
	Format(rec *Record) string
}

type DefaultFormatter struct {
	TimeFmt string
	Fmt     string
}

func NewDefaultFormatter(format string, timeFmt string) *DefaultFormatter {
	if format == "" {
		format = "%(msg)v"
	}
	if timeFmt == "" {
		timeFmt = time.Stamp
	}
	df := &DefaultFormatter{
		TimeFmt: timeFmt,
		Fmt:     format,
	}
	return df
}

var fieldHolderRegexp = regexp.MustCompile("\\$\\{\\w+\\}")

func (df *DefaultFormatter) Format(rec *Record) string {
	args := []interface{}{}
	newFmt := df.Fmt
	fieldMap := map[string]interface{}{
		"name":      rec.Name,
		"time":      rec.Time.Format(df.TimeFmt),
		"levelno":   rec.Level,
		"levelname": LevelMap[rec.Level],
		"lfile":     rec.LFile,
		"sfile":     rec.SFile,
		"line":      rec.Line,
		"msg":       rec.Message,
	}
	newFmt = strings.Replace(newFmt, "%", "%%", -1)
	newFmt = fieldHolderRegexp.ReplaceAllStringFunc(newFmt, func(match string) string {
		fieldName := match[2 : len(match)-1]
		field, ok := fieldMap[fieldName]
		if ok {
			args = append(args, field)
			return "%v"
		}
		return match
	})

	return fmt.Sprintf(newFmt, args...)
}

var defaultLevelColors = map[LogLevel]string{
	DebugLevel:    "bold_cyan",
	InfoLevel:     "bold_green",
	WarnLevel:     "bold_yellow",
	ErrorLevel:    "bold_red",
	CriticalLevel: "bg_red",
}

type RainbowFormatter struct {
	*DefaultFormatter
	LevelColors map[LogLevel]string
}

func NewRainbowFormatter(format string, timeFmt string) *RainbowFormatter {
	rf := &RainbowFormatter{
		DefaultFormatter: NewDefaultFormatter(format, timeFmt),
		LevelColors:      defaultLevelColors,
	}
	return rf
}

func (rf *RainbowFormatter) Format(rec *Record) string {
	newFmt := rf.DefaultFormatter.Format(rec)

	newFmt = fieldHolderRegexp.ReplaceAllStringFunc(newFmt, func(match string) string {
		m := match[2 : len(match)-1]
		if m == "log_color" {
			m, _ = rf.LevelColors[rec.Level]
		}
		code, ok := EscapeCodes[m]
		if ok {
			return code
		}
		return match
	})

	return newFmt
}
