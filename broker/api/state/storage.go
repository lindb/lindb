package state

import (
	"context"
	"net/http"

	"github.com/shirou/gopsutil/disk"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

// StorageAPI represents query storage cluster's state api from broker state machine
type StorageAPI struct {
	ctx                context.Context
	repo               state.Repository
	stateMachine       broker.StorageStateMachine
	shardAssignService service.ShardAssignService
	databaseService    service.DatabaseService
}

// NewStorageAPI creates storage state api
func NewStorageAPI(ctx context.Context, repo state.Repository,
	stateMachine broker.StorageStateMachine,
	shardAssignService service.ShardAssignService,
	databaseService service.DatabaseService,
) *StorageAPI {
	return &StorageAPI{
		ctx:                ctx,
		repo:               repo,
		stateMachine:       stateMachine,
		shardAssignService: shardAssignService,
		databaseService:    databaseService,
	}
}

// GetStorageClusterState returns the storage cluster detail stat by given cluster name
func (s *StorageAPI) GetStorageClusterState(w http.ResponseWriter, r *http.Request) {
	clusterName, err := api.GetParamsFromRequest("name", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	databaseList, shardAssignMap, err := s.getDatabaseInfo()
	if err != nil {
		api.Error(w, err)
		return
	}
	clusterStat, err := s.getStorageClusterInfo(clusterName)
	if err != nil {
		api.Error(w, err)
		return
	}
	aliveNodes := s.getStorageAliveNodes(clusterName)

	nodeStatMap := make(map[string]*models.NodeStat)
	for _, nodeStat := range clusterStat.Nodes {
		nodeStatMap[nodeStat.Node.Node.Indicator()] = nodeStat
	}

	for _, db := range databaseList {
		if db.Cluster != clusterName {
			continue
		}
		clusterStat.ReplicaStatus.Total += db.NumOfShard * db.ReplicaFactor

		shardAssign, ok := shardAssignMap[db.Name]
		if !ok {
			continue
		}
		db.Desc = db.String()
		databaseStatus := models.DatabaseStatus{
			Config:        *db,
			ReplicaStatus: models.ReplicaStatus{},
		}
		databaseStatus.ReplicaStatus.Total = db.NumOfShard * db.ReplicaFactor
		clusterStat.DatabaseStatusList = append(clusterStat.DatabaseStatusList, databaseStatus)

		shards := shardAssign.Shards
		nodes := shardAssign.Nodes
		for _, replica := range shards {
			available, underReplicated := calcReplicaStatus(replica, nodes, aliveNodes)
			for _, nodeID := range replica.Replicas {
				node := nodes[nodeID]
				nodeStat, ok := nodeStatMap[node.Indicator()]
				if !ok {
					continue
				}
				nodeStat.Replicas++
			}
			databaseStatus.ReplicaStatus.UnderReplicated += underReplicated
			clusterStat.ReplicaStatus.UnderReplicated += underReplicated
			if available == 0 {
				clusterStat.ReplicaStatus.Unavailable++
				databaseStatus.ReplicaStatus.Unavailable++
			}
		}
	}
	// calc node status
	clusterStat.NodeStatus.Total = len(clusterStat.Nodes)
	if aliveNodes != nil {
		clusterStat.NodeStatus.Alive = len(aliveNodes.ActiveNodes)
	}
	clusterStat.NodeStatus.Dead = clusterStat.NodeStatus.Total - clusterStat.NodeStatus.Alive

	api.OK(w, clusterStat)
}

// ListStorageCluster lists state of all storage clusters
func (s *StorageAPI) ListStorageClusterState(w http.ResponseWriter, r *http.Request) {
	databaseList, shardAssignMap, err := s.getDatabaseInfo()
	if err != nil {
		api.Error(w, err)
		return
	}
	clusterMap, err := s.getStorageClusterInfoMap()
	if err != nil {
		api.Error(w, err)
		return
	}
	storageNodeStatus := s.stateMachine.List()

	// calc node status
	aliveNodeMap := make(map[string]*models.StorageState)
	for _, storageState := range storageNodeStatus {
		aliveNodeMap[storageState.Name] = storageState
		clusterStat, ok := clusterMap[storageState.Name]
		if ok {
			clusterStat.NodeStatus.Alive = len(storageState.ActiveNodes)
			clusterStat.NodeStatus.Dead = clusterStat.NodeStatus.Total - clusterStat.NodeStatus.Alive
		}
	}

	for _, db := range databaseList {
		clusterStat, ok := clusterMap[db.Cluster]
		if !ok {
			continue
		}
		clusterStat.ReplicaStatus.Total += db.NumOfShard * db.ReplicaFactor

		shardAssign, ok := shardAssignMap[db.Name]
		if !ok {
			continue
		}
		aliveNodes, ok := aliveNodeMap[db.Cluster]
		if !ok {
			continue
		}
		shards := shardAssign.Shards
		nodes := shardAssign.Nodes
		for _, replica := range shards {
			available, underReplicated := calcReplicaStatus(replica, nodes, aliveNodes)
			clusterStat.ReplicaStatus.UnderReplicated += underReplicated
			if available == 0 {
				clusterStat.ReplicaStatus.Unavailable++
			}
		}
	}

	// build result
	var result []*models.StorageClusterStat
	for _, value := range clusterMap {
		result = append(result, value)
	}
	api.OK(w, result)
}

// getStorageAliveNodes returns the alive nodes of storage cluster by given cluster name
func (s *StorageAPI) getStorageAliveNodes(clusterName string) *models.StorageState {
	storageNodeStatus := s.stateMachine.List()
	var aliveNodes *models.StorageState
	for _, cluster := range storageNodeStatus {
		if cluster.Name == clusterName {
			aliveNodes = cluster
			break
		}
	}
	return aliveNodes
}

// getStorageClusterInfoMap returns the all storage cluster info
func (s *StorageAPI) getStorageClusterInfoMap() (clusterMap map[string]*models.StorageClusterStat, err error) {
	kvs, err := s.repo.List(s.ctx, constants.StorageClusterStatPath)
	if err != nil {
		return
	}

	clusterMap = make(map[string]*models.StorageClusterStat)

	for _, kv := range kvs {
		stat := models.StorageClusterStat{}
		err = encoding.JSONUnmarshal(kv.Value, &stat)
		if err != nil {
			return
		}
		diskUsageStat := disk.UsageStat{}
		for _, node := range stat.Nodes {
			diskUsageStat.Total += node.System.DiskUsageStat.Total
			diskUsageStat.Used += node.System.DiskUsageStat.Used
		}
		diskUsageStat.UsedPercent = float64(diskUsageStat.Used*100.0) / float64(diskUsageStat.Total)
		stat.Capacity = diskUsageStat

		nodeStatus := models.NodeStatus{}
		nodeStatus.Total = len(stat.Nodes)
		stat.NodeStatus = nodeStatus
		stat.ReplicaStatus = models.ReplicaStatus{}

		clusterMap[stat.Name] = &stat
	}
	return
}

// getStorageClusterInfo returns the storage cluster stat info by given cluster name
func (s *StorageAPI) getStorageClusterInfo(clusterName string) (stat *models.StorageClusterStat, err error) {
	statData, err := s.repo.Get(s.ctx, constants.GetStorageClusterStatPath(clusterName))
	if err != nil {
		return
	}
	stat = &models.StorageClusterStat{}
	err = encoding.JSONUnmarshal(statData, stat)
	if err != nil {
		return
	}
	diskUsageStat := disk.UsageStat{}
	for _, node := range stat.Nodes {
		diskUsageStat.Total += node.System.DiskUsageStat.Total
		diskUsageStat.Used += node.System.DiskUsageStat.Used
	}
	diskUsageStat.UsedPercent = float64(diskUsageStat.Used*100.0) / float64(diskUsageStat.Total)
	stat.Capacity = diskUsageStat
	return
}

// getDatabaseInfo returns the database info include database's config and shard assignment
func (s *StorageAPI) getDatabaseInfo() (databaseList []*models.Database, shardAssignMap map[string]*models.ShardAssignment, err error) {
	databaseList, err = s.databaseService.List()
	if err != nil {
		return
	}
	shardAssignList, err := s.shardAssignService.List()
	if err != nil {
		return
	}
	shardAssignMap = make(map[string]*models.ShardAssignment)
	for _, shardAssign := range shardAssignList {
		shardAssignMap[shardAssign.Name] = shardAssign
	}
	return
}

// nodeIsAlive returns the node if alive
func nodeIsAlive(storageState *models.StorageState, nodeID string) bool {
	if storageState == nil {
		return false
	}
	_, ok := storageState.ActiveNodes[nodeID]
	return ok
}

// calcReplicaStatus calculates the replica status
func calcReplicaStatus(replica *models.Replica,
	nodes map[int]*models.Node,
	storageState *models.StorageState,
) (available int, underReplicated int) {
	for _, nodeID := range replica.Replicas {
		node := nodes[nodeID]
		if nodeIsAlive(storageState, node.Indicator()) {
			available++
		} else {
			underReplicated++
		}
	}
	return
}
