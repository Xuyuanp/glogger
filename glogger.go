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
	"errors"
	"fmt"
	"sync"
)

type loggerMapper struct {
	mapper map[string]*Logger
	mu     sync.RWMutex
}

var lm *loggerMapper
var once sync.Once

func init() {
	once.Do(setup)
}

func setup() {
	lm = &loggerMapper{
		mapper: map[string]*Logger{},
	}
}

func (lm *loggerMapper) GetLogger(name string) *Logger {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.mapper[name]
}

func (lm *loggerMapper) registerLogger(logger *Logger) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	_, ok := lm.mapper[logger.Name]
	if ok {
		return errors.New(fmt.Sprintf("Logger with name %s has exists", logger.Name))
	}
	lm.mapper[logger.Name] = logger
	return nil
}

func GetLogger(name string) *Logger {
	return lm.GetLogger(name)
}

func registerLogger(logger *Logger) error {
	return lm.registerLogger(logger)
}
