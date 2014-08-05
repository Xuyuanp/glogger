package glogger

import (
	"io"
	"os"
	"sync"
)

type Handler interface {
	Handle(rec *Record)
}

type GenericHandler struct {
	Filterer
	Fmter Formatter
	mu    sync.Mutex
}

func (gh *GenericHandler) Format(rec *Record) string {
	return gh.Fmter.Format(rec)
}

func (gd *GenericHandler) Emit(text string) {
}

func (gh *GenericHandler) Handle(rec *Record) {
	if !gh.Filter(rec) {
		return
	}
	gh.mu.Lock()
	defer gh.mu.Unlock()
	text := gh.Format(rec)
	gh.Emit(text)
}

type StreamHandler struct {
	GenericHandler
	Writer io.Writer
}

func NewStreamHandler(w io.Writer) *StreamHandler {
	if w == nil {
		w = os.Stderr
	}
	sh := new(StreamHandler)
	sh.Writer = w
	sh.Fmter = NewDefaultFormatter("")
	return sh
}

func (sh *StreamHandler) Emit(text string) {
	sh.Writer.Write([]byte(text + "\n"))
}

func (sh *StreamHandler) Handle(rec *Record) {
	if !sh.Filter(rec) {
		return
	}
	sh.mu.Lock()
	defer sh.mu.Unlock()
	text := sh.Format(rec)
	sh.Emit(text)
}
