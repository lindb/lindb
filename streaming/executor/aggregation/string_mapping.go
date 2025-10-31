package aggregation

import (
	"strings"
	"sync"

	"go.uber.org/atomic"
)

type StringMapping struct {
	strToID map[string]uint32
	idToStr map[uint32]string

	id atomic.Uint32

	mutex sync.Mutex
}

func NewStringMapping() *StringMapping {
	return &StringMapping{
		strToID: make(map[string]uint32),
		idToStr: make(map[uint32]string),
	}
}

func (m *StringMapping) GetID(s string) uint32 {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if id, ok := m.strToID[s]; ok {
		return id
	}

	/*
		make a new copy for s in order to remove references from bigger string slice
		for example:
		big := "a very large string..."
		s := big[10:20]  // s is a substring (slice) of big
	*/
	newStr := strings.Clone(s)

	id := m.id.Inc()
	m.strToID[newStr] = id
	m.idToStr[id] = newStr
	return id
}

func (m *StringMapping) GetString(id uint32) string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.idToStr[id]
}
