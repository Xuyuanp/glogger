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
	"io"
	"os"
	"sync"
)

type Handler interface {
	Filter
	Emit(log string)
	Format(rec *Record) string
	GetLevel() LogLevel
	GetMutex() *sync.Mutex
	GetName() string
}

type HandlerGroup struct {
	Handlers *list.List
}

func (hg *HandlerGroup) AddHandler(h Handler) {
	if hg.Handlers == nil {
		hg.Handlers = list.New()
	}
	hg.Handlers.PushBack(h)
}

func (hg *HandlerGroup) Handle(rec *Record) {
	if hg.Handlers == nil {
		return
	}
	for e := hg.Handlers.Front(); e != nil; e = e.Next() {
		var h Handler = e.Value.(Handler)
		func() {
			if !h.DoFilter(rec) {
				return
			}
			h.GetMutex().Lock()
			defer h.GetMutex().Unlock()
			log := h.Format(rec)
			h.Emit(log)
		}()
	}
}

type GenericHandler struct {
	FilterGroup
	level     LogLevel
	name      string
	formatter Formatter
	mu        sync.Mutex
}

func NewHandler(name string, level LogLevel, formatter Formatter) *GenericHandler {
	gh := &GenericHandler{
		name:      name,
		level:     level,
		formatter: formatter,
	}
	gh.AddFilter(NewLevelFilter(level))
	return gh
}

func (gh *GenericHandler) Format(rec *Record) string {
	return gh.formatter.Format(rec)
}

func (gh *GenericHandler) GetLevel() LogLevel {
	return gh.level
}

func (gh *GenericHandler) GetMutex() *sync.Mutex {
	return &(gh.mu)
}

func (gh *GenericHandler) GetName() string {
	return gh.name
}

type StreamHandler struct {
	*GenericHandler
	Writer io.Writer
}

func NewStreamHandler(name string, level LogLevel, formatter Formatter, w io.Writer) *StreamHandler {
	if w == nil {
		panic(w)
	}
	sh := &StreamHandler{
		GenericHandler: NewHandler(name, level, formatter),
		Writer:         w,
	}
	return sh
}

func (sh *StreamHandler) Emit(text string) {
	sh.Writer.Write([]byte(text + "\n"))
}

type FileHandler struct {
	*StreamHandler
	FileName string
	Flag     int
	Pem      os.FileMode
}

func NewFileHandler(name string, level LogLevel, formatter Formatter, fileName string, flag int, pem os.FileMode) *FileHandler {
	file, err := os.OpenFile(fileName, flag, pem)
	if err != nil {
		panic(err)
	}
	fh := &FileHandler{
		StreamHandler: NewStreamHandler(name, level, formatter, file),
		FileName:      fileName,
		Flag:          flag,
		Pem:           pem,
	}
	return fh
}

func (fh *FileHandler) Emit(text string) {
	fh.StreamHandler.Emit(text)
}
