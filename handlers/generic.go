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

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/Xuyuanp/glogger"
)

// GenericHandler is an abstract struct which fully implemented Handler interface
// expected Emit method.
type GenericHandler struct {
	glogger.GroupFilter
	level     glogger.LogLevel
	name      string
	formatter glogger.Formatter
	mu        sync.Mutex
}

func NewHandler() *GenericHandler {
	gh := &GenericHandler{
		level: glogger.DebugLevel,
	}
	return gh
}

func (gh *GenericHandler) Format(rec *glogger.Record) string {
	return gh.formatter.Format(rec)
}

func (gh *GenericHandler) Name() string {
	return gh.name
}

func (gh *GenericHandler) Level() glogger.LogLevel {
	return gh.level
}

func (gh *GenericHandler) SetName(name string) {
	gh.name = name
}

func (gh *GenericHandler) SetLevel(level glogger.LogLevel) {
	gh.level = level
}

func (gh *GenericHandler) Mutex() *sync.Mutex {
	return &(gh.mu)
}

func (gh *GenericHandler) LoadConfig(config []byte) {
	var m map[string]interface{}
	err := json.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}
	gh.LoadConfigFromMap(m)
}

func (gh *GenericHandler) LoadConfigFromMap(config map[string]interface{}) {
	name, ok := config["name"]
	if ok {
		gh.name = name.(string)
	}
	level, ok := config["level"]
	if ok {
		gh.level = glogger.StringToLevel[level.(string)]
	}
	formatter, ok := config["formatter"]
	if ok {
		gh.formatter = glogger.GetFormatter(formatter.(string))
	}
}

func (gh *GenericHandler) LoadConfigFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	code, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	gh.LoadConfig(code)
}
