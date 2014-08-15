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
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/Xuyuanp/glogger"
)

func init() {
	glogger.RegisterConfigLoaderBuilder("github.com/Xuyuanp/glogger/handlers.StreamHandler", func() glogger.ConfigLoader {
		return NewStreamHandler()
	})
}

type StreamHandler struct {
	*GenericHandler
	Writer io.Writer
}

func NewStreamHandler() *StreamHandler {
	sh := &StreamHandler{
		GenericHandler: NewHandler(),
	}
	return sh
}

func (sh *StreamHandler) Emit(text string) {
	sh.Writer.Write([]byte(text + "\n"))
}

var writerMap = map[string]io.Writer{
	"stdout": os.Stdout,
	"stderr": os.Stderr,
}

func (sh *StreamHandler) LoadConfig(config []byte) {
	var m map[string]interface{}
	err := json.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}
	sh.LoadConfigFromMap(m)
}

func (sh *StreamHandler) LoadConfigFromMap(config map[string]interface{}) {
	sh.GenericHandler.LoadConfigFromMap(config)
	if writer, ok := config["writer"]; ok {
		if w, ok := writerMap[writer.(string)]; ok {
			sh.Writer = w
		} else {
			panic("unknown writer: " + writer.(string))
		}
	} else {
		panic("'writer' field is required")
	}
}

func (sh *StreamHandler) LoadConfigFromFile(fileName string) {
	if file, err := os.Open(fileName); err == nil {
		defer file.Close()
	} else {
		panic(err)
	}

	if code, err := ioutil.ReadAll(file); err == nil {
		l.LoadConfig(code)
	} else {
		panic(err)
	}
}
