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

func NewHandler(name string, level glogger.LogLevel, formatter glogger.Formatter) *GenericHandler {
	gh := &GenericHandler{
		name:      name,
		level:     level,
		formatter: formatter,
	}
	return gh
}

func (gh *GenericHandler) Format(rec *glogger.Record) string {
	return gh.formatter.Format(rec)
}

func (gh *GenericHandler) Level() glogger.LogLevel {
	return gh.level
}

func (gh *GenericHandler) Mutex() *sync.Mutex {
	return &(gh.mu)
}

func (gh *GenericHandler) Name() string {
	return gh.name
}
