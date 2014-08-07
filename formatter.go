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

import "fmt"

var LevelMap = map[LogLevel]string{
	DebugLevel:    "DBUG",
	InfoLevel:     "INFO",
	WarnLevel:     "WARN",
	CriticalLevel: "CRIT",
}

type Formatter interface {
	Format(rec *Record) string
}

type DefaultFormatter struct {
	Fmt string
}

func NewDefaultFormatter(format string) Formatter {
	if format == "" {
		format = "[%s\t%s\t%s\t%s\t:%d] %s"
	}
	df := &DefaultFormatter{
		Fmt: format,
	}
	return df
}

func (df *DefaultFormatter) Format(rec *Record) string {
	levelName, _ := LevelMap[rec.Level]
	return fmt.Sprintf(df.Fmt, rec.Name, rec.Time, levelName, rec.File, rec.Line, rec.Message)
}
