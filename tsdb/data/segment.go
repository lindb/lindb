package data

import (
	"github.com/eleme/lindb/storage"
	"fmt"
)

type segment struct {
	baseTime int64
	store    *storage.Store
}

func newSegment(baseTime int64) (*segment, error) {
	store, err := storage.NewStore("test", storage.StoreOption{Path: "../test_data"})
	if err != nil {
		return nil, fmt.Errorf("create storage error:%s", err)
	}
	return &segment{
		baseTime: baseTime,
		store:    store,
	}, nil
}

func (s *segment) getFamily() {

}
