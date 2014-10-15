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
	"regexp"
	"strings"
)

type Formatter interface {
	ConfigLoader
	Format(rec *Record) string
}

var formatterRegister = NewRegister()

func RegisterFormatter(name string, formatter Formatter) {
	formatterRegister.Register(name, formatter)
}

func GetFormatter(name string) Formatter {
	if v := formatterRegister.Get(name); v != nil {
		return v.(Formatter)
	}
	return nil
}

func init() {
	RegisterConfigLoaderBuilder("github.com/Xuyuanp/glogger.DefaultFormatter", func() ConfigLoader {
		return NewDefaultFormatter()
	})
}

type DefaultFormatter struct {
	TimeFmt string `json:timefmt`
	Fmt     string `json:fmt`
}

var DefaultFormat = "[${time} ${levelname} ${sfile}:${line} ${func}] ${msg}"
var DefaultTimeFormat = "2006-01-02 15:04:05"

func NewDefaultFormatter() *DefaultFormatter {
	df := &DefaultFormatter{
		TimeFmt: DefaultTimeFormat,
		Fmt:     DefaultFormat,
	}
	return df
}

var FieldHolderRegexp = regexp.MustCompile(`\$\{\w+\}`)

func (df *DefaultFormatter) Format(rec *Record) string {
	args := []interface{}{}
	newFmt := df.Fmt
	fieldMap := map[string]interface{}{
		"name":      rec.Name,
		"time":      rec.Time.Format(df.TimeFmt),
		"levelno":   rec.Level,
		"levelname": LevelToString[rec.Level],
		"lfile":     rec.LFile,
		"sfile":     rec.SFile,
		"func":      rec.Func,
		"line":      rec.Line,
		"msg":       rec.Message,
	}
	newFmt = strings.Replace(newFmt, "%", "%%", -1)
	newFmt = FieldHolderRegexp.ReplaceAllStringFunc(newFmt, func(match string) string {
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

func (df *DefaultFormatter) LoadConfigJson(config []byte) error {
	return json.Unmarshal(config, df)
}

func (df *DefaultFormatter) LoadConfig(config map[string]interface{}) error {
	if code, err := json.Marshal(config); err == nil {
		return df.LoadConfigJson(code)
	} else {
		return err
	}
	return nil
}
