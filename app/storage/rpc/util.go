package rpc

import (
	"errors"
	"fmt"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage"
	"github.com/lindb/lindb/storage/store"
)

// getOrCreateSegment returns write ahead log if it exists, else creates a new write ahead log.
func getOrCreateSegment(
	engine storage.Engine,
	database string,
	shardID models.ShardID,
	segmentTime int64,
	leader models.NodeID,
) (store.WriteAheadLog, error) {
	db, ok := engine.GetDatabase(database)
	if !ok {
		return nil, errors.New("database not exist")
	}

	shard, ok := db.GetShard(shardID)
	if !ok {
		fmt.Println(shardID)
		return nil, errors.New("shard not exist")
	}
	p, err := shard.GetOrCreatePartition(segmentTime)
	if err != nil {
		fmt.Printf("err0=%s\n", err)
		return nil, err
	}
	segment, err := p.GetOrCreateSegment(segmentTime)
	if err != nil {
		fmt.Printf("err1=%s\n", err)
		return nil, err
	}
	log, err := segment.GetOrCreateWAL(leader)
	if err != nil {
		fmt.Printf("err2=%s\n", err)
		return nil, err
	}
	return log, nil
}
