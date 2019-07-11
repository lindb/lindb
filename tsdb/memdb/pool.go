package memdb

import (
	"sync"
)

var (
	// tsStoresListPool is a set for storing []*fieldStore
	tsStoresListPool = _tsStoresListPool{pool: sync.Pool{}}
	// tsStoresListPool is a set for storing []*metricStore
	metricStoresListPool = _metricStoresListPool{pool: sync.Pool{}}
)

type _tsStoresListPool struct {
	pool sync.Pool
}

// get picks pointer to []*timeSeriesStore from the pool.
func (p *_tsStoresListPool) get(length int) *[]*timeSeriesStore {
	item := p.pool.Get()
	if item == nil {
		buf := make([]*timeSeriesStore, length)
		return &buf
	}
	buf := item.(*[]*timeSeriesStore)
	// cap is smaller than required size.
	if cap(*buf) < length {
		p.pool.Put(item)
		buf := make([]*timeSeriesStore, length)
		return &buf
	}
	*buf = (*buf)[:length]
	return buf
}

// put returns a tsStore list to the pool
func (p *_tsStoresListPool) put(buf *[]*timeSeriesStore) {
	*buf = (*buf)[:0]
	p.pool.Put(buf)
}

type _metricStoresListPool struct {
	pool sync.Pool
}

// get picks pointer to []*metricStore from the pool.
func (p *_metricStoresListPool) get(length int) *[]*metricStore {
	item := p.pool.Get()
	if item == nil {
		buf := make([]*metricStore, length)
		return &buf
	}
	buf := item.(*[]*metricStore)
	// cap is smaller than required size.
	if cap(*buf) < length {
		p.pool.Put(item)
		buf := make([]*metricStore, length)
		return &buf
	}
	*buf = (*buf)[:length]
	return buf
}

// put returns a metricStoreList to the pool
func (p *_metricStoresListPool) put(buf *[]*metricStore) {
	*buf = (*buf)[:0]
	p.pool.Put(buf)
}
