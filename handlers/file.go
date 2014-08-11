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
	"path/filepath"

	"github.com/Xuyuanp/glogger"
)

type FileHandler struct {
	*StreamHandler
	FileName string
}

func NewFileHandler(name string, level glogger.LogLevel, formatter glogger.Formatter, fileName string) *FileHandler {
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return nil
	}
	fh := &FileHandler{
		StreamHandler: NewStreamHandler(name, level, formatter, nil),
		FileName:      fileName,
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
