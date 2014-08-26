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
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/Xuyuanp/glogger"
)

func init() {
	glogger.RegisterConfigLoaderBuilder("github.com/Xuyuanp/glogger/handlers.StreamHandler", func() glogger.ConfigLoader {
		return NewStreamHandler()
	})
}

type StreamHandler struct {
	*GenericHandler
	nestedLogger *log.Logger
	mu           sync.Mutex
}

func NewStreamHandler() *StreamHandler {
	sh := &StreamHandler{
		GenericHandler: NewHandler(),
		nestedLogger:   log.New(os.Stdout, "", 0),
	}
	return sh
}

func (sh *StreamHandler) Handle(rec *glogger.Record) {
	sh.nestedLogger.Println(sh.Format(rec))
}

func (sh *StreamHandler) SetWriter(writer io.Writer) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.nestedLogger = log.New(writer, "", 0)
}

var writerMap = map[string]io.Writer{
	"stdout": os.Stdout,
	"stderr": os.Stderr,
}

func (sh *StreamHandler) LoadConfig(config map[string]interface{}) error {
	sh.GenericHandler.LoadConfig(config)
	if writer, ok := config["writer"]; ok {
		if w, ok := writerMap[writer.(string)]; ok {
			sh.SetWriter(w)
		} else {
			return fmt.Errorf("unknown writer: " + writer.(string))
		}
	} else {
		return fmt.Errorf("'writer' field is required")
	}
	return nil
}
