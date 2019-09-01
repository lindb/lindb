package tsdb

import "github.com/lindb/lindb/kv"

func init() {
	kv.RegisterMerger("mock_merger", &Merger{})
}

type Merger struct {
}

func (m *Merger) Merge(key uint32, value [][]byte) ([]byte, error) {
	//FIXME codingcrush
	return nil, nil
}
