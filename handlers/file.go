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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Xuyuanp/glogger"
)

func init() {
	glogger.RegisterConfigLoaderBuilder("github.com/Xuyuanp/glogger/handlers.FileHandler", func() glogger.ConfigLoader {
		return NewFileHandler()
	})
}

type FileHandler struct {
	*StreamHandler
	FileName string
}

func NewFileHandler() *FileHandler {
	fh := &FileHandler{
		StreamHandler: NewStreamHandler(),
	}
	return fh
}

func (fh *FileHandler) Emit(text string) {
	if fh.Writer == nil {
		file, err := os.OpenFile(fh.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fh.Writer = file
	}
	fh.StreamHandler.Emit(text)
}

func (fh *FileHandler) LoadConfig(config []byte) {
	var m map[string]interface{}
	if err := json.Unmarshal(config, &m); err == nil {
		fh.LoadConfigFromMap(m)
	} else {
		panic(err)
	}
}

func (fh *FileHandler) LoadConfigFromMap(config map[string]interface{}) {
	fh.GenericHandler.LoadConfigFromMap(config)
	if filename, ok := config["filename"]; ok {
		fh.FileName = filename.(string)
	} else {
		panic("'filename' field is required")
	}
}

func (fh *FileHandler) LoadConfigFromFile(fileName string) {
	if file, err := os.Open(fileName); err == nil {
		defer file.Close()
		if code, err := ioutil.ReadAll(file); err == nil {
			fh.LoadConfig(code)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
