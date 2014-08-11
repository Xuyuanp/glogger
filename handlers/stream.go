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
		GenericHandler: new(GenericHandler),
	}
	return sh
}

func (sh *StreamHandler) Emit(text string) {
	sh.Writer.Write([]byte(text + "\n"))
}

func (sh *StreamHandler) LoadConfig(config []byte) {
	var m map[string]interface{}
	err := json.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}
	sh.LoadConfigFromMap(m)
}

var writerMap = map[string]io.Writer{
	"stdout": os.Stdout,
	"stderr": os.Stderr,
}

func (sh *StreamHandler) LoadConfigFromMap(config map[string]interface{}) {
	sh.GenericHandler.LoadConfigFromMap(config)
	writer, ok := config["writer"]
	if ok {
		sh.Writer = writerMap[writer.(string)]
	}
}

func (sh *StreamHandler) LoadConfigFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	code, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	sh.LoadConfig(code)
}
