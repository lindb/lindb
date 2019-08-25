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
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/service"
)

//go:generate mockgen -source=./channel.go -destination=./channel_mock.go -package=replication

// ErrCanceled is the error returned when writing data ctx canceled.
var ErrCanceled = errors.New("write data ctx done")

const defaultReportInterval = 30

var log = logger.GetLogger("replication", "ChannelManager")

// ChannelManager manages the construction, retrieving, closing for all channels.
type ChannelManager interface {
	// GetChannel returns the channel for specific database and hash.
	// Error is returned when database is invalid or the total num of Channels
	// for a database is less than the numOfShard.
	GetChannel(database string, hash int32) (Channel, error)

	// CreateChannel creates a new channel or returns a existed channel for storage with specific database and shardID,
	// numOfShard should be greater or equal than the origin setting, otherwise error is returned.
	// numOfShard is used eot calculate the shardID for a given hash.
	CreateChannel(database string, numOfShard, shardID int32) (Channel, error)

	// Close closes all the channel.
	Close()
}

// channelManager implements ChannelManager.
type channelManager struct {
	// context passed to all Channel
	ctx context.Context
	// cancelFun to cancel context
	cancel context.CancelFunc
	// config
	cfg config.ReplicationChannel
	// factory to get rpc  write client
	fct rpc.ClientStreamFactory
	// for report replica state
	replicatorService service.ReplicatorService
	// channelID(a tuple of database, shardID)  -> Channel
	channelMap sync.Map
	// databaseID(a tuple of database)  -> numOfShard
	databaseShardsMap sync.Map
	// lock for channelMap
	lock4map sync.Mutex
}

// NewChannelManager returns a ChannelManager with dirPath and WriteClientFactory.
// WriteClientFactory makes it easy to mock rpc streamClient for test.
func NewChannelManager(cfg config.ReplicationChannel, fct rpc.ClientStreamFactory,
	replicatorService service.ReplicatorService) ChannelManager {
	ctx, cancel := context.WithCancel(context.Background())
	cm := &channelManager{
		ctx:               ctx,
		cancel:            cancel,
		cfg:               cfg,
		fct:               fct,
		replicatorService: replicatorService,
	}
	cm.scheduleStateReport()
	return cm
}

// GetChannel returns the channel for specific database and hash.
// Error is returned when database is invalid or the total num of Channels
// for a database is less than the numOfShard.
func (cm *channelManager) GetChannel(database string, hash int32) (Channel, error) {
	shardVal, ok := cm.databaseShardsMap.Load(database)
	if !ok {
		return nil, fmt.Errorf("channel for  database:%s not found", database)
	}
	numOfShard := shardVal.(int32)
	shardID := hash % numOfShard
	if shardID < 0 {
		shardID = -shardID
	}

	channelID := cm.buildChannelID(database, shardID)
	channelVal, ok := cm.channelMap.Load(channelID)

	if !ok {
		return nil, fmt.Errorf("channel for database:%s, shardID:%d not found", database, shardID)
	}

	ch := channelVal.(Channel)
	return ch, nil
}

// CreateChannel creates a new channel or returns a existed channel for storage with specific database and shardID.
// NumOfShard should be greater or equal than the origin setting, otherwise error is returned.
func (cm *channelManager) CreateChannel(database string, numOfShard, shardID int32) (Channel, error) {
	if numOfShard <= 0 || shardID >= numOfShard {
		return nil, errors.New("numOfShard should be greater than 0 and shardID should less then numOfShard")
	}
	channelID := cm.buildChannelID(database, shardID)
	val, ok := cm.channelMap.Load(channelID)
	if !ok {
		// double check
		cm.lock4map.Lock()
		defer cm.lock4map.Unlock()
		val, ok = cm.channelMap.Load(channelID)
		if !ok {
			// check numOfShard
			shardVal, ok := cm.databaseShardsMap.Load(database)
			if ok {
				oldNumOfShard := shardVal.(int32)
				if numOfShard < oldNumOfShard {
					return nil, errors.New("numOfShard should be equal or greater than original setting")
				}
			}
			cm.databaseShardsMap.Store(database, numOfShard)

			ch, err := newChannel(cm.ctx, cm.cfg, database, shardID, cm.fct)
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

// scheduleStateReport schedules a state report background job
func (cm *channelManager) scheduleStateReport() {
	interval := defaultReportInterval
	if cm.cfg.ReportInterval > 0 {
		interval = cm.cfg.ReportInterval
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				cm.reportState()
			case <-cm.ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// reportState reports the state of all replicators under current broker
func (cm *channelManager) reportState() {
	brokerState := models.BrokerReplicaState{
		ReportTime: timeutil.Now(),
	}
	cm.channelMap.Range(func(key, value interface{}) bool {
		channel, ok := value.(Channel)
		if ok {
			channel.Database()
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
				brokerState.Replicas = append(brokerState.Replicas, replicatorState)
			}
		}
		return true
	})
	if err := cm.replicatorService.Report(&brokerState); err != nil {
		log.Error("report broker replicator state fail", logger.Error(err))
	}
}

// buildChannelID return a string id by joining database, shardID with separator.
func (cm *channelManager) buildChannelID(database string, shardID int32) string {
	return database + "/" + strconv.Itoa(int(shardID))
}

// Channel represents a place to buffer the data for a specific cluster, database, shardID.
type Channel interface {
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
func newChannel(cxt context.Context, cfg config.ReplicationChannel, database string, shardID int32,
	fct rpc.ClientStreamFactory) (Channel, error) {
	dirPath := path.Join(cfg.Dir, database, strconv.Itoa(int(shardID)))
	interval := time.Duration(cfg.RemoveTaskIntervalInSecond) * time.Second

	q, err := queue.NewFanOutQueue(dirPath, cfg.SegmentFileSize, interval)
	if err != nil {
		return nil, err
	}

	c := &channel{
		cxt:      cxt,
		dirPath:  dirPath,
		fct:      fct,
		database: database,
		shardID:  shardID,
		q:        q,
		ch:       make(chan []byte, cfg.BufferSize),
		logger:   logger.GetLogger("replication", "Channel"),
	}

	c.initAppendTask()
	c.watchClose()

	return c, nil
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
			rep := newReplicator(target, c.database, c.shardID, fo, c.fct)

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
