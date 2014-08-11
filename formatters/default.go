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

package formatters

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Xuyuanp/glogger"
)

var LevelMap = map[glogger.LogLevel]string{
	glogger.DebugLevel:    "DBUG",
	glogger.InfoLevel:     "INFO",
	glogger.WarnLevel:     "WARN",
	glogger.ErrorLevel:    "ERRO",
	glogger.CriticalLevel: "CRIT",
}

type DefaultFormatter struct {
	TimeFmt string
	Fmt     string
}

var defaultFormat = "[${time} ${levelname} ${sfile}:${line} ${func}] ${msg}"
var defaultTimeFormat = "2006-01-02 15:04:05"

func NewDefaultFormatter() *DefaultFormatter {
	df := &DefaultFormatter{
		TimeFmt: defaultTimeFormat,
		Fmt:     defaultFormat,
	}
	return df
}

var fieldHolderRegexp = regexp.MustCompile("\\$\\{\\w+\\}")

func (df *DefaultFormatter) Format(rec *glogger.Record) string {
	args := []interface{}{}
	newFmt := df.Fmt
	fieldMap := map[string]interface{}{
		"name":      rec.Name,
		"time":      rec.Time.Format(df.TimeFmt),
		"levelno":   rec.Level,
		"levelname": LevelMap[rec.Level],
		"lfile":     rec.LFile,
		"sfile":     rec.SFile,
		"func":      rec.Func,
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
