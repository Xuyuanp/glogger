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

func init() {
	onceHandlerManager.Do(initHandlermanager)
}

type Handler interface {
	Namer
	Leveler
	Filter
	Emit(log string)
	Format(rec *Record) string
	Mutex() *sync.Mutex
}

type handlerManager struct {
	mapper map[string]Handler
	mu     sync.RWMutex
}

var hdlManager *handlerManager
var onceHandlerManager sync.Once

func initHandlermanager() {
	hdlManager = &handlerManager{
		mapper: map[string]Handler{},
	}
}

func RegisterHandler(name string, handler Handler) {
	hdlManager.mu.Lock()
	defer hdlManager.mu.Unlock()
	_, dup := hdlManager.mapper[name]
	if dup {
		panic("Register Handler named " + name + " twice")
	}
	hdlManager.mapper[name] = handler
}

func GetHandler(name string) Handler {
	hdlManager.mu.RLock()
	defer hdlManager.mu.RUnlock()
	return hdlManager.mapper[name]
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
		var h Handler = e.Value.(Handler)
		func() {
			if rec.Level < h.Level() || !h.Filter(rec) {
				return
			}
			h.Mutex().Lock()
			defer h.Mutex().Unlock()
			log := h.Format(rec)
			h.Emit(log)
		}()
	}
}
