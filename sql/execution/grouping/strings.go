package grouping

type StringMapper struct {
	s map[string]uint32
	i map[uint32]string

	idx uint32
}

func (sm *StringMapper) GetID(val string) uint32 {
	id, ok := sm.s[val]
	if ok {
		return id
	}

	id = sm.idx
	sm.idx++
	sm.s[val] = id
	sm.i[id] = val
	return id
}

func (sm *StringMapper) GetValue(id uint32) string {
	return sm.i[id]
}
