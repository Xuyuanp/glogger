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

import "github.com/Xuyuanp/glogger"

var defaultLevelColors = map[glogger.LogLevel]string{
	glogger.DebugLevel:    "bold_cyan",
	glogger.InfoLevel:     "bold_green",
	glogger.WarnLevel:     "bold_yellow",
	glogger.ErrorLevel:    "bold_red",
	glogger.CriticalLevel: "bg_bold_red",
}

var defaultRainbowFormat = "[${time} ${log_color}${levelname}${reset} ${dim}${green}${sfile}${reset}:${line} ${dim_cyan}${func}${reset}] ${msg}"

type RainbowFormatter struct {
	*DefaultFormatter
	LevelColors map[glogger.LogLevel]string
}

func NewRainbowFormatter() *RainbowFormatter {
	rf := &RainbowFormatter{
		DefaultFormatter: &DefaultFormatter{
			Fmt:     defaultRainbowFormat,
			TimeFmt: defaultTimeFormat,
		},
		LevelColors: defaultLevelColors,
	}
	return rf
}

func (rf *RainbowFormatter) Format(rec *glogger.Record) string {
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

	newFmt += EscapeCodes["reset"]

	return newFmt
}