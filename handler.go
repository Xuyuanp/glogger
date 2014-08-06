package glogger

import (
	"container/list"
	"io"
	"os"
	"sync"
)

type Handler interface {
	GetMutex() *sync.Mutex
	Emit(log string)
	SetFormatter(fmt Formatter)
	Format(rec *Record) string
}

type HandlerGroup struct {
	FilterGroup
	Level    LogLevel
	Handlers *list.List
}

func (hg *HandlerGroup) AddHandler(h Handler) {
	if hg.Handlers == nil {
		hg.Handlers = list.New()
	}
	hg.Handlers.PushBack(h)
}

func (hg *HandlerGroup) Handle(rec *Record) {
	if rec.Level < hg.Level || !hg.DoFilter(rec) || hg.Handlers == nil {
		return
	}
	for e := hg.Handlers.Front(); e != nil; e = e.Next() {
		var h Handler = e.Value.(Handler)
		func() {
			h.GetMutex().Lock()
			defer h.GetMutex().Unlock()
			log := h.Format(rec)
			h.Emit(log)
		}()
	}
}

type GenericHandler struct {
	Fmter Formatter
	mu    sync.Mutex
}

func (gh *GenericHandler) Format(rec *Record) string {
	return gh.Fmter.Format(rec)
}

func (gh *GenericHandler) GetMutex() *sync.Mutex {
	return &(gh.mu)
}

func (gh *GenericHandler) SetFormatter(fmt Formatter) {
	gh.Fmter = fmt
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

type FileHandler struct {
	StreamHandler
	FileName string
	Flag     int
	Pem      os.FileMode
}

func NewFileHandler(fileName string, flag int, pem os.FileMode) *FileHandler {
	fh := &FileHandler{
		FileName: fileName,
		Flag:     flag,
		Pem:      pem,
	}
	fh.Fmter = NewDefaultFormatter("")
	return fh
}

func (fh *FileHandler) Emit(text string) {
	if fh.Writer == nil {
		file, _ := os.OpenFile(fh.FileName, fh.Flag, fh.Pem)
		fh.Writer = file
	}
	fh.StreamHandler.Emit(text)
}
