package replication

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./channel.go -destination=./channel_mock.go -package=replication

// ErrCanceled is the error returned when writing data ctx canceled.
var ErrCanceled = errors.New("write data ctx done")

// ChannelManager manages the construction, retrieving, closing for all channels.
type ChannelManager interface {
	// GetChannel returns the channel for specific cluster, database and hash.
	// Error is returned when cluster, database is invalid or the total num of Channels
	// for a cluster, database is less than the numOfShard.
	GetChannel(cluster, database string, hash int32) (Channel, error)

	// CreateChannel creates a new channel or returns a existed channel for storage with specific cluster,
	// database and shardID, numOfShard should be greater or equal than the origin setting, otherwise error is returned.
	// numOfShard is used eot calculate the shardID for a given hash.
	CreateChannel(cluster, database string, numOfShard, shardID int32) (Channel, error)

	// Close closes all the channel.
	Close()
}

// channelManager implements ChannelManager.
type channelManager struct {
	// context passed to all Channel
	cxt context.Context
	// cancelFun to cancel context
	cancel context.CancelFunc
	// config
	cfg config.ReplicationChannel
	// factory to get rpc  write client
	fct rpc.ClientStreamFactory
	// channelID(a tuple of cluster, database, shardID)  -> Channel
	channelMap sync.Map
	// databaseID(a tuple of cluster, database)  -> numOfShard
	databaseShardsMap sync.Map
	// lock for channelMap
	lock4map sync.Mutex
}

// NewChannelManager returns a ChannelManager with dirPath and WriteClientFactory.
// WriteClientFactory makes it easy to mock rpc streamClient for test.
func NewChannelManager(cfg config.ReplicationChannel, fct rpc.ClientStreamFactory) ChannelManager {
	cxt, cancel := context.WithCancel(context.TODO())
	return &channelManager{
		cxt:    cxt,
		cancel: cancel,
		cfg:    cfg,
		fct:    fct,
	}
}

// GetChannel returns the channel for specific cluster, database and hash.
// Error is returned when cluster, database is invalid or the total num of Channels
// for a cluster, database is less than the numOfShard.
func (cm *channelManager) GetChannel(cluster, database string, hash int32) (Channel, error) {
	shardVal, ok := cm.databaseShardsMap.Load(cm.buildDatabaseID(cluster, database))
	if !ok {
		return nil, fmt.Errorf("channel for cluster:%s, database:%s not found", cluster, database)
	}
	numOfShard := shardVal.(int32)
	shardID := hash % numOfShard
	if shardID < 0 {
		shardID = -shardID
	}

	channelID := cm.buildChannelID(cluster, database, shardID)
	channelVal, ok := cm.channelMap.Load(channelID)

	if !ok {
		return nil, fmt.Errorf("channel for cluster:%s, database:%s, shardID:%d not found", cluster, database, shardID)
	}

	ch := channelVal.(Channel)
	return ch, nil
}

// CreateChannel creates a new channel or returns a existed channel for storage with specific cluster,
// database and shardID. NumOfShard should be greater or equal than the origin setting, otherwise error is returned.
func (cm *channelManager) CreateChannel(cluster, database string, numOfShard, shardID int32) (Channel, error) {
	if numOfShard <= 0 || shardID >= numOfShard {
		return nil, errors.New("numOfShard should be greater than 0 and shardID should less then numOfShard")
	}
	channelID := cm.buildChannelID(cluster, database, shardID)
	val, ok := cm.channelMap.Load(channelID)
	if !ok {
		// double check
		cm.lock4map.Lock()
		defer cm.lock4map.Unlock()
		val, ok = cm.channelMap.Load(channelID)
		if !ok {
			// check numOfShard
			dbID := cm.buildDatabaseID(cluster, database)
			shardVal, ok := cm.databaseShardsMap.Load(dbID)
			if ok {
				oldNumOfShard := shardVal.(int32)
				if numOfShard < oldNumOfShard {
					return nil, errors.New("numOfShard should be equal or greater than original setting")
				}
			}
			cm.databaseShardsMap.Store(dbID, numOfShard)

			ch, err := newChannel(cm.cxt, cm.cfg,
				cluster, database, shardID, cm.fct)
			if err != nil {
				return nil, err
			}
			cm.channelMap.Store(channelID, ch)
			return ch, nil
		}
	}

	ch := val.(Channel)
	return ch, nil
}

// Close closes all the channel.
func (cm *channelManager) Close() {
	cm.cancel()
}

// buildChannelID return a string id by joining cluster, database, shardID with separator.
func (cm *channelManager) buildChannelID(cluster, database string, shardID int32) string {
	return cluster + "/" + database + "/" + strconv.Itoa(int(shardID))
}

// buildDatabaseID return a string id by joining cluster, database.
func (cm *channelManager) buildDatabaseID(cluster, database string) string {
	return cluster + "/" + database
}

// Channel represents a place to buffer the data for a specific cluster, database, shardID.
type Channel interface {
	// Cluster returns the cluster attribution.
	Cluster() string
	// Database returns the database attribution.
	Database() string
	// ShardID returns the shardID attribution.
	ShardID() int32
	// Write writes the data into the channel, ErrCanceled is returned when the ctx is canceled before
	// data is wrote successfully.
	// Concurrent safe.
	Write(cxt context.Context, data []byte) error
	// GetOrCreateReplicator get a existed or creates a new replicator for target.
	// Concurrent safe.
	GetOrCreateReplicator(target models.Node) (Replicator, error)
	// Nodes returns all the target nodes for replication.
	Targets() []models.Node
}

// channel implements Channel.
type channel struct {
	// context to close channel
	cxt     context.Context
	dirPath string
	// factory to get WriteClient
	fct      rpc.ClientStreamFactory
	cluster  string
	database string
	shardID  int32
	// underlying storage for written data
	q queue.FanOutQueue
	// chanel to convert multiple goroutine write to single goroutine write to FanOutQueue
	ch chan []byte
	// target -> replicator map
	replicatorMap sync.Map
	// lock to protect replicatorMap
	lock4map sync.RWMutex
	logger   *logger.Logger
}

// newChannel returns a new channel with specific attribution.
func newChannel(cxt context.Context, cfg config.ReplicationChannel, cluster, database string, shardID int32,
	fct rpc.ClientStreamFactory) (Channel, error) {
	dirPath := path.Join(cfg.Path, cluster, database, strconv.Itoa(int(shardID)))
	interval := time.Duration(cfg.RemoveTaskIntervalInSecond) * time.Second

	q, err := queue.NewFanOutQueue(dirPath, cfg.SegmentFileSize, interval)
	if err != nil {
		return nil, err
	}

	c := &channel{
		cxt:      cxt,
		dirPath:  dirPath,
		fct:      fct,
		cluster:  cluster,
		database: database,
		shardID:  shardID,
		q:        q,
		ch:       make(chan []byte, cfg.BufferSize),
		logger:   logger.GetLogger("replication/channel"),
	}

	c.initAppendTask()
	c.watchClose()

	return c, nil
}

// Cluster returns the cluster attribution.
func (c *channel) Cluster() string {
	return c.cluster
}

// Database returns the database attribution.
func (c *channel) Database() string {
	return c.database
}

// ShardID returns the shardID attribution.
func (c *channel) ShardID() int32 {
	return c.shardID
}

// GetOrCreateReplicator get a existed or creates a new replicator for target.
// Concurrent safe.
func (c *channel) GetOrCreateReplicator(target models.Node) (Replicator, error) {
	val, ok := c.replicatorMap.Load(target)
	if !ok {
		// double check
		c.lock4map.Lock()
		defer c.lock4map.Unlock()
		val, ok = c.replicatorMap.Load(target)
		if !ok {
			fo, err := c.q.GetOrCreateFanOut(target.Indicator())
			if err != nil {
				return nil, err
			}
			rep := newReplicator(target, c.cluster, c.database, c.shardID, fo, c.fct)

			c.replicatorMap.Store(target, rep)
			return rep, nil
		}
	}
	rep := val.(Replicator)
	return rep, nil
}

// Nodes returns all the nodes for replication.
func (c *channel) Targets() []models.Node {
	nodes := make([]models.Node, 0)
	c.replicatorMap.Range(func(key, value interface{}) bool {
		nd, _ := key.(models.Node)
		nodes = append(nodes, nd)
		return true
	})
	return nodes
}

// Write writes the data into the channel, ErrCanceled is returned when the ctx is canceled before
// data is wrote successfully.
// Concurrent safe.
func (c *channel) Write(cxt context.Context, data []byte) error {
	select {
	case c.ch <- data:
		return nil
	case <-cxt.Done():
		return ErrCanceled
	}
}

// initAppendTask starts a goroutine to consume data from ch and append to q.
func (c *channel) initAppendTask() {
	go func() {
		for data := range c.ch {
			_, err := c.q.Append(data)
			if err != nil {
				// todo retry?
				c.logger.Error("append data error", logger.String("dirPath", c.dirPath),
					logger.Error(err))
			}
		}

		c.logger.Info("close queue for channel", logger.String("dirPath", c.dirPath))
		c.q.Close()
	}()
}

// watchClose waits on the context done then close the ch.
func (c *channel) watchClose() {
	go func() {
		<-c.cxt.Done()
		c.lock4map.RLock()
		defer c.lock4map.RUnlock()
		c.replicatorMap.Range(func(key, value interface{}) bool {
			rep, _ := value.(Replicator)
			rep.Stop()
			return true
		})
		// todo avoid Write send data to closed channel.
		close(c.ch)
	}()
}
