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

import "container/list"

type Filter interface {
	DoFilter(rec *Record) bool
}

type FilterGroup struct {
	filters *list.List
}

func (f *FilterGroup) AddFilter(ft Filter) {
	if f.filters == nil {
		f.filters = list.New()
	}
	f.filters.PushBack(ft)
}

func (f *FilterGroup) DoFilter(rec *Record) bool {
	if f.filters == nil {
		return true
	}
	for e := f.filters.Front(); e != nil; e = e.Next() {
		filter := e.Value.(Filter)
		if !filter.DoFilter(rec) {
			return false
		}
	}
	return true
}

type LevelFilter struct {
	Level LogLevel
}

func NewLevelFilter(level LogLevel) *LevelFilter {
	return &LevelFilter{
		Level: level,
	}
}

func (lf *LevelFilter) DoFilter(rec *Record) bool {
	return lf.Level <= rec.Level
}
