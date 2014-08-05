package glogger

type Filterer struct {
	children []*Filterer
}

func (f *Filterer) AddChild(child *Filterer) {
	if f.children == nil {
		f.children = make([]*Filterer, 5)
	}
	f.children = append(f.children, child)
}

func (f *Filterer) RemoveChild(child *Filterer) {
	if f.children == nil {
		f.children = make([]*Filterer, 5)
	}
	for i, ch := range f.children {
		if ch == child {
			f.children = append(f.children[:i], f.children[i+1:]...)
		}
	}
}

func (f *Filterer) Filter(rec *Record) bool {
	for _, child := range f.children {
		if !child.Filter(rec) {
			return false
		}
	}
	return true
}
