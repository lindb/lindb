// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package state

//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/shirou/gopsutil/disk"
//
//	"github.com/lindb/lindb/app/broker/api/admin"
//	"github.com/lindb/lindb/app/broker/deps"
//	"github.com/lindb/lindb/constants"
//	"github.com/lindb/lindb/models"
//	"github.com/lindb/lindb/pkg/encoding"
//	"github.com/lindb/lindb/pkg/http"
//	"github.com/lindb/lindb/pkg/logger"
//)
//
//var (
//	StorageStatePath     = "/storage/cluster/state"
//	ListStorageStatePath = "/storage/cluster/state/list"
//)
//
//// StorageAPI represents query storage cluster's state api from broker state machine
//type StorageAPI struct {
//	deps         *deps.HTTPDeps
//	ListDatabase func() ([]*models.Database, error)
//	logger       *logger.Logger
//}
//
//// NewStorageAPI creates storage state api
//func NewStorageAPI(deps *deps.HTTPDeps) *StorageAPI {
//	return &StorageAPI{
//		deps:         deps,
//		ListDatabase: admin.NewDatabaseAPI(deps).ListDataBase,
//		logger:       logger.GetLogger("broker", "StorageAPI"),
//	}
//}
//
//// Register adds storage state url route.
//func (s *StorageAPI) Register(route gin.IRoutes) {
//	route.GET(StorageStatePath, s.GetStorageClusterState)
//	route.GET(ListStorageStatePath, s.ListStorageClusterState)
//}
//
//// GetStorageClusterState returns the storage cluster detail stat by given cluster name
//func (s *StorageAPI) GetStorageClusterState(c *gin.Context) {
//	var param struct {
//		ClusterName string `form:"name" binding:"required"`
//	}
//	err := c.ShouldBindQuery(&param)
//	if err != nil {
//		http.Error(c, err)
//		return
//	}
//	databaseList, shardAssignMap, err := s.getDatabaseInfo()
//	if err != nil {
//		http.Error(c, err)
//		return
//	}
//	clusterStat, err := s.getStorageClusterInfo(param.ClusterName)
//	if err != nil {
//		http.Error(c, err)
//		return
//	}
//	aliveNodes := s.getStorageAliveNodes(param.ClusterName)
//
//	nodeStatMap := make(map[string]*models.NodeStat)
//	for _, nodeStat := range clusterStat.Nodes {
//		nodeStatMap[nodeStat.Node.Node.Indicator()] = nodeStat
//	}
//
//	for _, db := range databaseList {
//		if db.Cluster != param.ClusterName {
//			continue
//		}
//		clusterStat.ReplicaStatus.Total += db.NumOfShard * db.ReplicaFactor
//
//		shardAssign, ok := shardAssignMap[db.Name]
//		if !ok {
//			continue
//		}
//		db.Desc = db.String()
//		databaseStatus := models.DatabaseStatus{
//			Config:        *db,
//			ReplicaStatus: models.ReplicaStatus{},
//		}
//		databaseStatus.ReplicaStatus.Total = db.NumOfShard * db.ReplicaFactor
//		clusterStat.DatabaseStatusList = append(clusterStat.DatabaseStatusList, databaseStatus)
//
//		shards := shardAssign.Shards
//		nodes := shardAssign.Nodes
//		for _, replica := range shards {
//			available, underReplicated := calcReplicaStatus(replica, nodes, aliveNodes)
//			for _, nodeID := range replica.Replicas {
//				node := nodes[nodeID]
//				nodeStat, ok := nodeStatMap[node.Indicator()]
//				if !ok {
//					continue
//				}
//				nodeStat.Replicas++
//			}
//			databaseStatus.ReplicaStatus.UnderReplicated += underReplicated
//			clusterStat.ReplicaStatus.UnderReplicated += underReplicated
//			if available == 0 {
//				clusterStat.ReplicaStatus.Unavailable++
//				databaseStatus.ReplicaStatus.Unavailable++
//			}
//		}
//	}
//	// calc node status
//	clusterStat.NodeStatus.Total = len(clusterStat.Nodes)
//	if aliveNodes != nil {
//		clusterStat.NodeStatus.Alive = len(aliveNodes.ActiveNodes)
//	}
//	clusterStat.NodeStatus.Dead = clusterStat.NodeStatus.Total - clusterStat.NodeStatus.Alive
//
//	http.OK(c, clusterStat)
//}
//
//// ListStorageClusterState lists state of all storage clusters
//func (s *StorageAPI) ListStorageClusterState(c *gin.Context) {
//	databaseList, shardAssignMap, err := s.getDatabaseInfo()
//	if err != nil {
//		http.Error(c, err)
//		return
//	}
//	clusterMap, err := s.getStorageClusterInfoMap()
//	if err != nil {
//		http.Error(c, err)
//		return
//	}
//	storageNodeStatus := s.deps.StateMachines.StorageSM.List()
//
//	// calc node status
//	aliveNodeMap := make(map[string]*models.StorageState)
//	for _, storageState := range storageNodeStatus {
//		aliveNodeMap[storageState.Name] = storageState
//		clusterStat, ok := clusterMap[storageState.Name]
//		if ok {
//			clusterStat.NodeStatus.Alive = len(storageState.ActiveNodes)
//			clusterStat.NodeStatus.Dead = clusterStat.NodeStatus.Total - clusterStat.NodeStatus.Alive
//		}
//	}
//
//	for _, db := range databaseList {
//		clusterStat, ok := clusterMap[db.Cluster]
//		if !ok {
//			continue
//		}
//		clusterStat.ReplicaStatus.Total += db.NumOfShard * db.ReplicaFactor
//
//		shardAssign, ok := shardAssignMap[db.Name]
//		if !ok {
//			continue
//		}
//		aliveNodes, ok := aliveNodeMap[db.Cluster]
//		if !ok {
//			continue
//		}
//		shards := shardAssign.Shards
//		nodes := shardAssign.Nodes
//		for _, replica := range shards {
//			available, underReplicated := calcReplicaStatus(replica, nodes, aliveNodes)
//			clusterStat.ReplicaStatus.UnderReplicated += underReplicated
//			if available == 0 {
//				clusterStat.ReplicaStatus.Unavailable++
//			}
//		}
//	}
//
//	// build result
//	var result []*models.StorageClusterStat
//	for _, value := range clusterMap {
//		result = append(result, value)
//	}
//	http.OK(c, result)
//}
//
//// getStorageAliveNodes returns the alive nodes of storage cluster by given cluster name
//func (s *StorageAPI) getStorageAliveNodes(clusterName string) *models.StorageState {
//	storageNodeStatus := s.deps.StateMachines.StorageSM.List()
//	var aliveNodes *models.StorageState
//	for _, cluster := range storageNodeStatus {
//		if cluster.Name == clusterName {
//			aliveNodes = cluster
//			break
//		}
//	}
//	return aliveNodes
//}
//
//// getStorageClusterInfoMap returns the all storage cluster info
//func (s *StorageAPI) getStorageClusterInfoMap() (clusterMap map[string]*models.StorageClusterStat, err error) {
//	ctx, cancel := s.deps.WithTimeout()
//	defer cancel()
//	kvs, err := s.deps.Repo.List(ctx, constants.StorageClusterStatPath)
//	if err != nil {
//		return
//	}
//
//	clusterMap = make(map[string]*models.StorageClusterStat)
//
//	for _, kv := range kvs {
//		stat := models.StorageClusterStat{}
//		err = encoding.JSONUnmarshal(kv.Value, &stat)
//		if err != nil {
//			return
//		}
//		diskUsageStat := disk.UsageStat{}
//		for _, node := range stat.Nodes {
//			diskUsageStat.Total += node.System.DiskUsageStat.Total
//			diskUsageStat.Used += node.System.DiskUsageStat.Used
//		}
//		diskUsageStat.UsedPercent = float64(diskUsageStat.Used*100.0) / float64(diskUsageStat.Total)
//		stat.Capacity = diskUsageStat
//
//		nodeStatus := models.NodeStatus{}
//		nodeStatus.Total = len(stat.Nodes)
//		stat.NodeStatus = nodeStatus
//		stat.ReplicaStatus = models.ReplicaStatus{}
//
//		clusterMap[stat.Name] = &stat
//	}
//	return clusterMap, err
//}
//
//// getStorageClusterInfo returns the storage cluster stat info by given cluster name
//func (s *StorageAPI) getStorageClusterInfo(clusterName string) (stat *models.StorageClusterStat, err error) {
//	ctx, cancel := s.deps.WithTimeout()
//	defer cancel()
//	statData, err := s.deps.Repo.Get(ctx, constants.GetStorageClusterStatPath(clusterName))
//	if err != nil {
//		return
//	}
//	stat = &models.StorageClusterStat{}
//	err = encoding.JSONUnmarshal(statData, stat)
//	if err != nil {
//		return
//	}
//	diskUsageStat := disk.UsageStat{}
//	for _, node := range stat.Nodes {
//		diskUsageStat.Total += node.System.DiskUsageStat.Total
//		diskUsageStat.Used += node.System.DiskUsageStat.Used
//	}
//	diskUsageStat.UsedPercent = float64(diskUsageStat.Used*100.0) / float64(diskUsageStat.Total)
//	stat.Capacity = diskUsageStat
//	return
//}
//
//// getDatabaseInfo returns the database info include database's config and shard assignment
//func (s *StorageAPI) getDatabaseInfo() (
//	databaseList []*models.Database,
//	shardAssignMap map[string]*models.ShardAssignment,
//	err error,
//) {
//	ctx, cancel := s.deps.WithTimeout()
//	defer cancel()
//
//	data, err := s.deps.Repo.List(ctx, constants.DatabaseAssignPath)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	var result []*models.ShardAssignment
//	for _, val := range data {
//		shardAssign := &models.ShardAssignment{}
//		err = encoding.JSONUnmarshal(val.Value, shardAssign)
//		if err != nil {
//			s.logger.Warn("unmarshal data error",
//				logger.String("data", string(val.Value)))
//		} else {
//			result = append(result, shardAssign)
//		}
//	}
//
//	databaseList, err = s.ListDatabase()
//	if err != nil {
//		return
//	}
//
//	shardAssignMap = make(map[string]*models.ShardAssignment)
//	for _, shardAssign := range result {
//		shardAssignMap[shardAssign.Name] = shardAssign
//	}
//	return databaseList, shardAssignMap, err
//}
//
//// nodeIsAlive returns the node if alive
//func nodeIsAlive(storageState *models.StorageState, nodeID string) bool {
//	if storageState == nil {
//		return false
//	}
//	_, ok := storageState.ActiveNodes[nodeID]
//	return ok
//}
//
//// calcReplicaStatus calculates the replica status
//func calcReplicaStatus(replica *models.Replica,
//	nodes map[models.NodeID]*models.Node,
//	storageState *models.StorageState,
//) (available int, underReplicated int) {
//	for _, nodeID := range replica.Replicas {
//		node := nodes[nodeID]
//		if nodeIsAlive(storageState, node.Indicator()) {
//			available++
//		} else {
//			underReplicated++
//		}
//	}
//	return
//}
