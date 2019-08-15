package replica

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/pathutil"
)

//go:generate mockgen -source=./status_state_machine.go -destination=./status_state_machine_mock.go -package=replica

// StatusStateMachine represents the status of database's replicas
// Each broker node need start this state machine,
type StatusStateMachine interface {
	discovery.Listener
	// GetQueryableReplicas returns the queryable replicas，
	// and chooses the fastest replica if the shard has multi-replica.
	// returns storage node => shard id list
	GetQueryableReplicas(database string) map[string][]int32
	// GetReplicas returns the replica state list under this broker by broker's indicator
	GetReplicas(broker string) models.BrokerReplicaState
	// Close closes state machine, stops watch change event
	Close() error
}

// statusStateMachine implements status state machine,
// watches replica state path for listening modify event which broker uploaded
type statusStateMachine struct {
	discovery discovery.Discovery

	ctx    context.Context
	cancel context.CancelFunc

	mutex sync.RWMutex
	// brokers: broker node => replica list under this broker
	brokers map[string]models.BrokerReplicaState

	log *logger.Logger
}

// NewStatusStateMachine creates a replica's status state machine
func NewStatusStateMachine(ctx context.Context, factory discovery.Factory) (StatusStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	sm := &statusStateMachine{
		ctx:     c,
		cancel:  cancel,
		brokers: make(map[string]models.BrokerReplicaState),
		log:     logger.GetLogger("replica/status/state/machine"),
	}
	// new replica status discovery
	sm.discovery = factory.CreateDiscovery(constants.ReplicaStatePath, sm)
	if err := sm.discovery.Discovery(); err != nil {
		return nil, fmt.Errorf("discovery database status error:%s", err)
	}
	return sm, nil
}

// GetQueryableReplicas returns the queryable replicas
// returns storage node => shard id list
func (sm *statusStateMachine) GetQueryableReplicas(database string) map[string][]int32 {
	// 1. find shards by given database's name
	shards := make(map[string][]models.ReplicaState)
	sm.mutex.RLock()
	for _, brokerReplicaState := range sm.brokers {
		for _, replica := range brokerReplicaState.Replicas {
			if replica.Database != database {
				continue
			}
			shardID := replica.ShardIndicator()
			shards[shardID] = append(shards[shardID], replica)
		}
	}
	sm.mutex.RUnlock()

	if len(shards) == 0 {
		return nil
	}

	result := make(map[string][]int32)
	for _, replicas := range shards {
		replicaList := replicas
		if len(replicaList) > 1 {
			// has multi-replica, chooses the fastest
			// sort replicas based pending msg
			sort.Slice(replicaList, func(i, j int) bool {
				return replicaList[i].Pending < replicaList[j].Pending
			})
		}
		nodeID := replicaList[0].Target.Indicator()
		result[nodeID] = append(result[nodeID], replicaList[0].ShardID)
	}

	return result
}

// GetReplicas returns the replica state list under this broker by broker's indicator
func (sm *statusStateMachine) GetReplicas(broker string) models.BrokerReplicaState {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.brokers[broker]
}

// Close closes state machine, stops watch change event
func (sm *statusStateMachine) Close() error {
	sm.discovery.Close()
	sm.cancel()
	return nil
}

// OnCreates updates the broker's replica status when broker upload replica state
func (sm *statusStateMachine) OnCreate(key string, resource []byte) {
	brokerReplicaState := models.BrokerReplicaState{}
	if err := json.Unmarshal(resource, &brokerReplicaState); err != nil {
		sm.log.Error("discovery replica status but unmarshal error",
			logger.String("data", string(resource)), logger.Error(err))
		return
	}
	broker := pathutil.GetName(key)
	sm.mutex.Lock()
	sm.brokers[broker] = brokerReplicaState
	sm.mutex.Unlock()
}

// OnDelete deletes the broker's replica status when broker offline
func (sm *statusStateMachine) OnDelete(key string) {
	broker := pathutil.GetName(key)
	sm.mutex.Lock()
	delete(sm.brokers, broker)
	sm.mutex.Unlock()
}
