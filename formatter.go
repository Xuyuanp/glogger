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

import "sync"

func init() {
	onceFormatterManager.Do(initFormatterManager)
}

type Formatter interface {
	ConfigLoader
	Format(rec *Record) string
}

type formatterManager struct {
	mapper map[string]Formatter
	mu     sync.RWMutex
}

var fmtManager *formatterManager
var onceFormatterManager sync.Once

func initFormatterManager() {
	fmtManager = &formatterManager{
		mapper: map[string]Formatter{},
	}
}

func RegisterFormatter(name string, formatter Formatter) {
	fmtManager.mu.Lock()
	defer fmtManager.mu.Unlock()
	_, dup := fmtManager.mapper[name]
	if dup {
		panic("Formatter named " + name + " twice")
	}
	fmtManager.mapper[name] = formatter
}

func GetFormatter(name string) Formatter {
	fmtManager.mu.RLock()
	defer fmtManager.mu.RUnlock()
	return fmtManager.mapper[name]
}
