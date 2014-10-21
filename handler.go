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
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

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

func (hg *handlerGroup) SetHandlers(handlers ...Handler) {
	hg.ClearHandlers()
	for _, h := range handlers {
		hg.handlers.PushBack(h)
	}
}

func (hg *handlerGroup) ClearHandlers() {
	hg.handlers = list.New()
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

// GenericHandler is an abstract struct which fully implemented Handler interface
// expected Emit method.
type GenericHandler struct {
	GroupFilter
	level     LogLevel
	formatter Formatter
}

// NewHandler return a new GenericHandler
func NewHandler() *GenericHandler {
	gh := &GenericHandler{
		level:     DebugLevel,
		formatter: NewDefaultFormatter(),
	}
	return gh
}

// Format a record with formatter
func (gh *GenericHandler) Format(rec *Record) string {
	return gh.formatter.Format(rec)
}

// SetFormatter set a new Formatter
func (gh *GenericHandler) SetFormatter(formatter Formatter) {
	gh.formatter = formatter
}

// Level return log level of the handler
func (gh *GenericHandler) Level() LogLevel {
	return gh.level
}

// SetLevel set log level of handler
func (gh *GenericHandler) SetLevel(level LogLevel) {
	gh.level = level
}

// LoadConfig load configuration from a map
func (gh *GenericHandler) LoadConfig(config map[string]interface{}) error {
	// Load log level, default DebugLevel
	if level, ok := config["level"]; ok {
		if l, ok := StringToLevel[level.(string)]; ok {
			gh.level = l
		} else {
			return fmt.Errorf("unknown log level: %s", level.(string))
		}
	} else {
		gh.level = DebugLevel
	}
	// Load Formatter, default is DefaultFormatter
	if formatter, ok := config["formatter"]; ok {
		if f := GetFormatter(formatter.(string)); f != nil {
			gh.formatter = f
		} else {
			return fmt.Errorf("unknown formater name: " + formatter.(string))
		}
	} else {
		gh.formatter = NewDefaultFormatter()
	}
	// Load filters
	if filters, ok := config["filters"]; ok {
		if len(filters.([]interface{})) > 0 {
			for _, fname := range filters.([]interface{}) {
				if filter := GetFilter(fname.(string)); filter != nil {
					gh.AddFilter(filter)
				} else {
					return fmt.Errorf("unknown filter name: %s", fname.(string))
				}
			}
		}
	}
	return nil
}

func init() {
	RegisterConfigLoaderBuilder("github.com/Xuyuanp/glogger.StreamHandler", func() ConfigLoader {
		return NewStreamHandler()
	})
}

// StreamHandler struct
type StreamHandler struct {
	*GenericHandler
	nestedLogger *log.Logger
	mu           sync.Mutex
}

// NewStreamHandler return a new StreamHandler
func NewStreamHandler() *StreamHandler {
	sh := &StreamHandler{
		GenericHandler: NewHandler(),
		nestedLogger:   log.New(os.Stdout, "", 0),
	}
	return sh
}

// Handle a Record
func (sh *StreamHandler) Handle(rec *Record) {
	sh.nestedLogger.Println(sh.Format(rec))
}

// SetWriter set a output writer
func (sh *StreamHandler) SetWriter(writer io.Writer) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.nestedLogger = log.New(writer, "", 0)
}

var writerMap = map[string]io.Writer{
	"stdout": os.Stdout,
	"stderr": os.Stderr,
}

// LoadConfig load configuration from a map
func (sh *StreamHandler) LoadConfig(config map[string]interface{}) error {
	sh.GenericHandler.LoadConfig(config)
	if writer, ok := config["writer"]; ok {
		if w, ok := writerMap[writer.(string)]; ok {
			sh.SetWriter(w)
		} else {
			return fmt.Errorf("unknown writer: " + writer.(string))
		}
	} else {
		sh.SetWriter(os.Stdout)
	}
	return nil
}
