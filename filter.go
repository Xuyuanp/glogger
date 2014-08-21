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
	ConfigLoader
	Filter(rec *Record) bool
}

var filterRegister = NewRegister()

func RegisterFilter(name string, filter Filter) {
	filterRegister.Register(name, filter)
}

func GetFilter(name string) Filter {
	if v := filterRegister.Get(name); v != nil {
		return v.(Filter)
	}
	return nil
}

type GroupFilter struct {
	filters *list.List
}

func (f *GroupFilter) AddFilter(ft Filter) {
	if f.filters == nil {
		f.filters = list.New()
	}
	f.filters.PushBack(ft)
}

func (f *GroupFilter) Filter(rec *Record) bool {
	if f.filters == nil {
		return true
	}
	for e := f.filters.Front(); e != nil; e = e.Next() {
		filter := e.Value.(Filter)
		if !filter.Filter(rec) {
			return false
		}
	}
	return true
}
func (f *GroupFilter) LoadConfig(config map[string]interface{}) {

}
