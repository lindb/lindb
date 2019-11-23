package roaring

type shortIterable interface {
	hasNext() bool
	next() uint16
}

type shortPeekable interface {
	shortIterable
	peekNext() uint16
	advanceIfNeeded(minval uint16)
}

type shortIterator struct {
	slice []uint16
	loc   int
}

func (si *shortIterator) hasNext() bool {
	return si.loc < len(si.slice)
}

func (si *shortIterator) next() uint16 {
	a := si.slice[si.loc]
	si.loc++
	return a
}

func (si *shortIterator) peekNext() uint16 {
	return si.slice[si.loc]
}

func (si *shortIterator) advanceIfNeeded(minval uint16) {
	if si.hasNext() && si.peekNext() < minval {
		si.loc = advanceUntil(si.slice, si.loc, len(si.slice), minval)
	}
}

type reverseIterator struct {
	slice []uint16
	loc   int
}

func (si *reverseIterator) hasNext() bool {
	return si.loc >= 0
}

func (si *reverseIterator) next() uint16 {
	a := si.slice[si.loc]
	si.loc--
	return a
}

/**********************************************************************/
// customize for LinDB
/**********************************************************************/

// PeekableShortIterator represents the wrapper interface of shortPeekable
type PeekableShortIterator interface {
	// HasNext returns if has next element
	HasNext() bool
	// Next returns the next element
	Next() uint16
	// PeekNext peeks the next element
	PeekNext() uint16
	// AdvanceIfNeeded skips to min value
	AdvanceIfNeeded(minVal uint16)
}

type peekableShortIterator struct {
	it shortPeekable
}

func newPeekableShortIterator(it shortPeekable) PeekableShortIterator {
	return &peekableShortIterator{it: it}
}

func (p *peekableShortIterator) HasNext() bool {
	return p.it.hasNext()
}

func (p *peekableShortIterator) Next() uint16 {
	return p.it.next()
}

func (p *peekableShortIterator) PeekNext() uint16 {
	return p.it.peekNext()
}

func (p *peekableShortIterator) AdvanceIfNeeded(minVal uint16) {
	p.it.advanceIfNeeded(minVal)
}
