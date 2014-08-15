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
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Xuyuanp/glogger"
)

func init() {
	glogger.RegisterConfigLoaderBuilder("github.com/Xuyuanp/glogger/formatters.RainbowFormatter", func() glogger.ConfigLoader {
		return NewRainbowFormatter()
	})
}

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

func (rf *RainbowFormatter) LoadConfig(config []byte) {
	var m map[string]interface{}
	err := json.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}
	rf.LoadConfigFromMap(m)
}

func (rf *RainbowFormatter) LoadConfigFromMap(config map[string]interface{}) {
	if fmt, ok := config["fmt"]; ok {
		rf.Fmt = fmt.(string)
	}
	if timefmt, ok := config["timefmt"]; ok {
		rf.TimeFmt = timefmt.(string)
	}
	if colors, ok := config["colors"]; ok {
		colorConfig := colors.(map[string]interface{})
		for name, level := range glogger.StringToLevel {
			if colori, yes := colorConfig[name]; yes {
				rf.LevelColors[level] = colori.(string)
			} else {
				panic("unknown color: " + name)
			}
		}
	}
}

func (rf *RainbowFormatter) LoadConfigFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	code, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	rf.LoadConfig(code)
}
