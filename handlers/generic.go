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
	"fmt"

	"github.com/Xuyuanp/glogger"
)

// GenericHandler is an abstract struct which fully implemented Handler interface
// expected Emit method.
type GenericHandler struct {
	glogger.GroupFilter
	level     glogger.LogLevel
	formatter glogger.Formatter
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

func (gh *GenericHandler) SetFormatter(formatter glogger.Formatter) {
	gh.formatter = formatter
}

func (gh *GenericHandler) Level() glogger.LogLevel {
	return gh.level
}

func (gh *GenericHandler) SetLevel(level glogger.LogLevel) {
	gh.level = level
}

func (gh *GenericHandler) LoadConfig(config map[string]interface{}) error {
	if level, ok := config["level"]; ok {
		if l, ok := glogger.StringToLevel[level.(string)]; ok {
			gh.level = l
		} else {
			return fmt.Errorf("unknown level: " + level.(string))
		}
	} else {
		return fmt.Errorf("'level' field is required")
	}
	if formatter, ok := config["formatter"]; ok {
		if f := glogger.GetFormatter(formatter.(string)); f != nil {
			gh.formatter = f
		} else {
			return fmt.Errorf("unknown formater name: " + formatter.(string))
		}
	} else {
		return fmt.Errorf("'formater' field is required")
	}
	return nil
}
