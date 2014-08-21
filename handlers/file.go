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

func (fh *FileHandler) LoadConfig(config map[string]interface{}) {
	fh.GenericHandler.LoadConfig(config)
	if filename, ok := config["filename"]; ok {
		fh.FileName = filename.(string)
	} else {
		panic("'filename' field is required")
	}
}
