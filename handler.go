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

// Handler determines where the log message to output
type Handler interface {
	Leveler
	Filter
	Handle(rec *Record)
}

var handlerRegister = NewRegister()

// RegisterHandler register a Handler to global manager with specific name.
// The Handler registered can be accessed by GetHandler method anywhere with this name.
func RegisterHandler(name string, handler Handler) {
	handlerRegister.Register(name, handler)
}

// GetHandler return the Handler registered with this name.
// nil will by returned if no Handler registered with this name.
func GetHandler(name string) Handler {
	return handlerRegister.Get(name).(Handler)
	if v := handlerRegister.Get(name); v != nil {
		return v.(Handler)
	}
	return nil
}

type handlerGroup struct {
	handlers *list.List
}

func (hg *handlerGroup) AddHandler(h Handler) {
	if hg.handlers == nil {
		hg.handlers = list.New()
	}
	hg.handlers.PushBack(h)
}

func (hg *handlerGroup) Handle(rec *Record) {
	if hg.handlers == nil {
		return
	}
	for e := hg.handlers.Front(); e != nil; e = e.Next() {
		var h = e.Value.(Handler)
		func() {
			if rec.Level < h.Level() || !h.Filter(rec) {
				return
			}
			h.Handle(rec)
		}()
	}
}
