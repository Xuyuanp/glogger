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
	"container/list"

	"sync"
)

type Handler interface {
	Filter
	Emit(log string)
	Format(rec *Record) string
	Level() LogLevel
	Mutex() *sync.Mutex
	Name() string
}

type handlerManger struct {
	handlers *list.List
}

func (hm *handlerManger) AddHandler(h Handler) {
	if hm.handlers == nil {
		hm.handlers = list.New()
	}
	hm.handlers.PushBack(h)
}

func (hm *handlerManger) Handle(rec *Record) {
	if hm.handlers == nil {
		return
	}
	for e := hm.handlers.Front(); e != nil; e = e.Next() {
		var h Handler = e.Value.(Handler)
		func() {
			if rec.Level < h.Level() || !h.DoFilter(rec) {
				return
			}
			h.Mutex().Lock()
			defer h.Mutex().Unlock()
			log := h.Format(rec)
			h.Emit(log)
		}()
	}
}
