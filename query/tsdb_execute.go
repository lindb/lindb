package query

import (
	"fmt"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/tsdb"
)

// tsdbExecute represents execution search logic in tsdb level,
// does query task async, then merge result, such as map-reduce job
type tsdbExecute struct {
	engine   tsdb.Engine
	query    models.Query
	shardIDs []int32

	shards []tsdb.Shard

	err error
}

// NewTSDBExecutor creates execution which queries tsdb storage
func NewTSDBExecutor(engine tsdb.Engine, shardIDs []int32, query models.Query) Executor {
	return &tsdbExecute{
		engine:   engine,
		shardIDs: shardIDs,
		query:    query,
	}
}

// Execute executes search logic in tsdb level,
// 1) valition input params
// 2) build execute plan
// 3) build execute pipeline
// 4) run pipeline
func (e *tsdbExecute) Execute() {
	// do query validation
	if err := e.validation(); err != nil {
		e.err = err
		return
	}

	// get shard by given query shard id list
	for _, shardID := range e.shardIDs {
		shard := e.engine.GetShard(shardID)
		// if shard exist, add shard to query list
		if shard != nil {
			e.shards = append(e.shards, shard)
		}
	}

	// check got shards if valid
	if err := e.checkShards(); err != nil {
		e.err = err
		return
	}
}

// validation validates query input params and tsdb data are valid
func (e *tsdbExecute) validation() error {
	// check input shardIDs if empty
	if len(e.shardIDs) == 0 {
		return fmt.Errorf("there is no shard id in search condition")
	}
	// check engine has shard
	if e.engine.NumOfShards() == 0 {
		return fmt.Errorf("tsdb engine[%s] hasn't shard", e.engine.Name())
	}
	return nil
}

// checkShards checks got shards if valid
func (e *tsdbExecute) checkShards() error {
	numOfShards := len(e.shards)
	numOfShardIDs := len(e.shardIDs)
	if numOfShards == 0 {
		return fmt.Errorf("cannot find shard by given shard id")
	}
	if numOfShards != numOfShardIDs {
		return fmt.Errorf("got shard size[%d] not eqauls input shard id size[%d]", numOfShards, numOfShardIDs)
	}
	return nil
}
