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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Xuyuanp/glogger"
)

func init() {
	glogger.RegisterConfigLoaderBuilder("github.com/Xuyuanp/glogger/handlers.RotatingFileHandler", func() glogger.ConfigLoader {
		return NewRotatingFileHandler()
	})
}

// RotatingFileHandler struct
type RotatingFileHandler struct {
	*glogger.GenericHandler
	FileName       string
	File           *os.File
	AutoRotate     bool
	MaxSize        uint64
	MaxLine        uint64
	Daily          bool
	nextRotateTime time.Time
	currentSize    uint64
	currentLine    uint64
	BackupCount    int
	mu             sync.Mutex
}

// NewRotatingFileHandler return a new RotatingFileHandler
func NewRotatingFileHandler() *RotatingFileHandler {
	fh := &RotatingFileHandler{
		GenericHandler: glogger.NewHandler(),
		AutoRotate:     true,
		Daily:          true,
	}
	return fh
}

// Handle a record
func (fh *RotatingFileHandler) Handle(rec *glogger.Record) {
	fh.mu.Lock()
	defer fh.mu.Unlock()
	if fh.File == nil {
		fmt.Fprintln(os.Stderr, "No log file")
		return
	}
	msg := fh.Format(rec)
	fh.currentLine += uint64(len(strings.Split(msg, "\n")))
	fh.currentSize += uint64(len(msg))
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fh.File.WriteString(msg)

	if fh.checkRotate() {
		fh.doRotate()
	}
}

// SetFileName set the name of file to output
func (fh *RotatingFileHandler) setFileName(fileName string) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if fh.AutoRotate {
		fh.currentSize = 0
		fh.currentLine = 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fh.currentSize += uint64(len(line))
			fh.currentLine++
		}
		fh.setupNextRotateTime()
	}

	if fh.File != nil {
		fh.File.Close()
	}

	fh.File = file
	fh.FileName = fileName
}

func (fh *RotatingFileHandler) checkRotate() bool {
	if !fh.AutoRotate {
		return false
	}
	if fh.MaxLine > 0 && fh.currentLine >= fh.MaxLine {
		return true
	}
	if fh.MaxSize > 0 && fh.currentSize >= fh.MaxSize {
		return true
	}
	if fh.Daily {
		now := time.Now()
		if now.After(fh.nextRotateTime) {
			return true
		}
	}
	return false
}

func (fh *RotatingFileHandler) doRotate() {
	nextFileName := ""
	for i := 1; fh.BackupCount == 0 || i <= fh.BackupCount; i++ {
		fileName := fmt.Sprintf("%s.%d", fh.FileName, i)
		if _, err := os.Open(fileName); err != nil {
			nextFileName = fileName
			break
		}
	}
	if nextFileName == "" {
		fmt.Fprintln(os.Stderr, "No more backup file")
		fh.File.Close()
		fh.File = nil
		return
	}
	fh.File.Close()
	fh.File = nil
	os.Rename(fh.FileName, nextFileName)

	fh.setFileName(fh.FileName)
	fh.currentLine = 0
	fh.currentSize = 0
	fh.setupNextRotateTime()
}

func (fh *RotatingFileHandler) setupNextRotateTime() {
	now := time.Now()
	nextTimeStr := fmt.Sprintf("%d-%d-%d 00:00:00", now.Year(), now.Month(), now.Day())
	nextTime, _ := time.ParseInLocation("2006-1-2 15:04:05", nextTimeStr, time.Local)
	nextTime = nextTime.Add(24 * time.Hour)
	fh.nextRotateTime = nextTime
}

// LoadConfig load configuration from a map
func (fh *RotatingFileHandler) LoadConfig(config map[string]interface{}) error {
	if err := fh.GenericHandler.LoadConfig(config); err != nil {
		return err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, fh); err != nil {
		return err
	}
	if fh.FileName == "" {
		return fmt.Errorf("'filename' field is required")
	}
	fh.setFileName(fh.FileName)
	return nil
}
