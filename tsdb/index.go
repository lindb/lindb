package tsdb

import (
	"fmt"
	"path/filepath"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/tsdb/metadb"
)

//go:generate mockgen -source=./index.go -destination=./index_mock.go -package=tsdb

// Index represents an index include id sequencer and index creation for shard level
type Index interface {
	// GetIDSequencer returns id sequencer for metric level
	GetIDSequencer() metadb.IDSequencer
	// CreateIndexDatabase creates index database for shard level
	CreateIndexDatabase(shardID int32) (metadb.IndexDatabase, error)
	// Close closes index kv store
	Close()
}

// index implements Index interface, using common kv store for index storage
type index struct {
	indexStore kv.Store
	sequencer  metadb.IDSequencer
}

// newIndex creates an index
func newIndex(name string, cfg config.Engine) (Index, error) {
	storeOption := kv.DefaultStoreOption(filepath.Join(cfg.Dir, name, "index"))
	indexStore, err := kv.NewStore(name, storeOption)
	if err != nil {
		return nil, err
	}
	familyOption := kv.FamilyOption{
		CompactThreshold: 0,
		Merger:           "mock_merger", //FIXME codingcrush
	}
	metricMetaFamily, err := indexStore.CreateFamily("metric_meta", familyOption)
	if err != nil {
		return nil, err
	}
	metricIDsFamily, err := indexStore.CreateFamily("metric_ids", familyOption)
	if err != nil {
		return nil, err
	}

	return &index{
		indexStore: indexStore,
		sequencer:  metadb.NewIDSequencer(metricIDsFamily, metricMetaFamily),
	}, err
}

// GetIDSequencer returns id sequencer for metric level
func (i *index) GetIDSequencer() metadb.IDSequencer {
	return i.sequencer
}

// CreateIndexDatabase creates index database for shard level
func (i *index) CreateIndexDatabase(shardID int32) (metadb.IndexDatabase, error) {
	familyOption := kv.FamilyOption{
		CompactThreshold: 0,
		Merger:           "mock_merger", //FIXME codingcrush
	}
	invertedFamily, err := i.indexStore.CreateFamily(fmt.Sprintf("inverted_%d", shardID), familyOption)
	if err != nil {
		return nil, err
	}
	forwardFamily, err := i.indexStore.CreateFamily(fmt.Sprintf("forward_%d", shardID), familyOption)
	if err != nil {
		return nil, err
	}
	return metadb.NewIndexDatabase(i.sequencer, invertedFamily, forwardFamily), nil
}

// Close closes index kv store
func (i *index) Close() {
	if err := i.indexStore.Close(); err != nil {
		log.Error("close index kv store error", logger.Error(err))
	}
}
