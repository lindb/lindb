package replication

import (
	"context"
	"errors"
	"path"
	"sync"

	"github.com/cespare/xxhash"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source=./database_channel.go -destination=./database_channel_mock.go -package=replication

var (
	// define error types
	errChannelNotFound = errors.New("shard replica channel not found")
	errInvalidShardID  = errors.New("numOfShard should be greater than 0 and shardID should less then numOfShard")
	errInvalidShardNum = errors.New("numOfShard should be equal or greater than original setting")
)

// for testing
var (
	mkdir         = fileutil.MkDirIfNotExist
	createChannel = newChannel
)

// DatabaseChannel represents the database level replication channel
type DatabaseChannel interface {
	// Write writes the metric data into channel's buffer
	Write(metricList *field.MetricList) error
	// CreateChannel creates the shard level replication channel by given shard id
	CreateChannel(numOfShard, shardID int32) (Channel, error)
	// ReplicaState returns the replica state
	ReplicaState() (replicas []models.ReplicaState)
}

type databaseChannel struct {
	database      string
	ctx           context.Context
	cfg           config.ReplicationChannel
	fct           rpc.ClientStreamFactory
	numOfShard    atomic.Int32
	shardChannels sync.Map
	mutex         sync.Mutex
}

// newDatabaseChannel creates a new database replication channel
func newDatabaseChannel(ctx context.Context,
	database string, cfg config.ReplicationChannel, numOfShard int32,
	fct rpc.ClientStreamFactory,
) (DatabaseChannel, error) {
	dirPath := path.Join(cfg.Dir, database)
	if err := mkdir(dirPath); err != nil {
		return nil, err
	}
	ch := &databaseChannel{
		database: database,
		ctx:      ctx,
		cfg:      cfg,
		fct:      fct,
	}
	ch.numOfShard.Store(numOfShard)
	return ch, nil
}

// Write writes the metric data into channel's buffer
func (dc *databaseChannel) Write(metricList *field.MetricList) (err error) {
	// sharding metrics to shards
	numOfShard := uint64(dc.numOfShard.Load())
	for _, metric := range metricList.Metrics {
		hash := xxhash.Sum64String(tag.Concat(metric.Tags))
		// set tags hash code for storage side reuse
		// !!!IMPORTANT: storage side will use this hash for write
		metric.TagsHash = hash
		shardID := int32(hash % numOfShard)
		channel, ok := dc.getChannelByShardID(shardID)
		if !ok {
			err = errChannelNotFound
			// broker error, do not return to client
			log.Error("channel not found", logger.String("database", dc.database), logger.Int32("shardID", shardID))
			continue
		}
		if err = channel.Write(metric); err != nil {
			log.Error("channel write data error", logger.String("database", dc.database), logger.Int32("shardID", shardID))
		}
	}
	return
}

// CreateChannel creates the shard level replication channel by given shard id
func (dc *databaseChannel) CreateChannel(numOfShard, shardID int32) (Channel, error) {
	channel, ok := dc.getChannelByShardID(shardID)
	if !ok {
		dc.mutex.Lock()
		defer dc.mutex.Unlock()

		// double check
		channel, ok = dc.getChannelByShardID(shardID)
		if !ok {
			if numOfShard <= 0 || shardID >= numOfShard {
				return nil, errInvalidShardID
			}
			if numOfShard < dc.numOfShard.Load() {
				return nil, errInvalidShardNum
			}
			ch, err := createChannel(dc.ctx, dc.cfg, dc.database, shardID, dc.fct)
			if err != nil {
				return nil, err
			}
			// need startup channel
			ch.Startup()
			// cache shard level channel
			dc.shardChannels.Store(shardID, ch)
			return ch, nil
		}
	}
	return channel, nil
}

// ReplicaState returns the replica state
func (dc *databaseChannel) ReplicaState() (replicas []models.ReplicaState) {
	dc.shardChannels.Range(func(key, value interface{}) bool {
		channel, ok := value.(Channel)
		if ok {
			targets := channel.Targets()
			for i := range targets {
				target := targets[i]
				replicator, err := channel.GetOrCreateReplicator(target)
				if err != nil {
					log.Error("get replicator fail", logger.String("target", (&target).Indicator()), logger.Error(err))
					continue
				}
				replicatorState := models.ReplicaState{
					Database:     replicator.Database(),
					Target:       target,
					ShardID:      replicator.ShardID(),
					Pending:      replicator.Pending(),
					ReplicaIndex: replicator.ReplicaIndex(),
					AckIndex:     replicator.AckIndex(),
				}
				replicas = append(replicas, replicatorState)
			}
		}
		return true
	})
	return
}

// getChannelByShardID gets the replica channel by shard id
func (dc *databaseChannel) getChannelByShardID(shardID int32) (Channel, bool) {
	channel, ok := dc.shardChannels.Load(shardID)
	if !ok {
		return nil, ok
	}
	ch, ok := channel.(Channel)
	if !ok {
		return nil, ok
	}
	return ch, true
}
