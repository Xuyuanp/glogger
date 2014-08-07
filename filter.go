package glogger

import "container/list"

type Filter interface {
	DoFilter(rec *Record) bool
}

type FilterGroup struct {
	filters *list.List
}

func (f *FilterGroup) AddFilter(ft Filter) {
	if f.filters == nil {
		f.filters = list.New()
	}
	f.filters.PushBack(ft)
}

func (f *FilterGroup) DoFilter(rec *Record) bool {
	if f.filters == nil {
		return true
	}
	for e := f.filters.Front(); e != nil; e = e.Next() {
		filter := e.Value.(Filter)
		if !filter.DoFilter(rec) {
			return false
		}
	}
	return true
}

type LevelFilter struct {
	Level LogLevel
}

func NewLevelFilter(level LogLevel) *LevelFilter {
	return &LevelFilter{
		Level: level,
	}
}

func (lf *LevelFilter) DoFilter(rec *Record) bool {
	return lf.Level < rec.Level
}
