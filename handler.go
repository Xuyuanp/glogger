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
	"encoding/base64"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"strings"
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

type GroupHandler struct {
	Handlers *list.List
}

func (gh *GroupHandler) AddHandler(h Handler) {
	if gh.Handlers == nil {
		gh.Handlers = list.New()
	}
	gh.Handlers.PushBack(h)
}

func (gh *GroupHandler) Handle(rec *Record) {
	if gh.Handlers == nil {
		return
	}
	for e := gh.Handlers.Front(); e != nil; e = e.Next() {
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

// GenericHandler is an abstract struct which fully implemented Handler interface
// expected Emit method.
type GenericHandler struct {
	GroupFilter
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

type SmtpHandler struct {
	*GenericHandler
	Host    string
	Port    int
	From    string
	To      []string
	Auth    smtp.Auth
	Subject string
}

func NewSmtpHandler(name string, level LogLevel, formatter Formatter, host string, port int, from string, to []string, auth smtp.Auth, subject string) *SmtpHandler {
	sh := &SmtpHandler{
		GenericHandler: NewHandler(name, level, formatter),
		Host:           host,
		Port:           port,
		From:           from,
		To:             to,
		Auth:           auth,
		Subject:        subject,
	}
	return sh
}

func (sh *SmtpHandler) Emit(text string) {
	header := make(map[string]string)
	header["From"] = sh.From
	header["To"] = strings.Join(sh.To, ";")
	header["Subject"] = sh.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\t\n", k, v)
	}

	message += "\t\n" + base64.StdEncoding.EncodeToString([]byte(text))

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", sh.Host, sh.Port),
		sh.Auth,
		sh.From,
		sh.To,
		[]byte(message),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
