package field

import (
	"bytes"
	"sort"
)

type Field struct {
	Name  []byte
	Type  Type
	Value interface{}
}

// Fields implements sort.Interface
type Fields []Field

func (fs Fields) Len() int           { return len(fs) }
func (fs Fields) Swap(i, j int)      { fs[i], fs[j] = fs[j], fs[i] }
func (fs Fields) Less(i, j int) bool { return bytes.Compare(fs[i].Name, fs[j].Name) < 0 }

func (fs Fields) Search(name []byte) (idx int, ok bool) {
	idx = sort.Search(fs.Len(), func(i int) bool {
		return bytes.Compare(fs[i].Name, name) >= 0
	})
	if idx >= fs.Len() || !bytes.Equal(fs[idx].Name, name) {
		return -1, false
	}
	return idx, true
}

// Insert adds or replace a Field
func (fs Fields) Insert(f Field) Fields {
	idx, ok := fs.Search(f.Name)
	if !ok {
		next := append(fs, f)
		sort.Sort(next)
		return next
	}
	fs[idx] = f
	return fs
}
