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

package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc/proto/common"
	replicaRpc "github.com/lindb/lindb/rpc/proto/replica"
	"github.com/lindb/lindb/rpc/proto/storage"
)

//go:generate mockgen -source ./rpc.go -destination=./rpc_mock.go -package=rpc

const (
	metaKeyLogicNode = "metaKeyLogicNode"
	metaKeyDatabase  = "metaKeyDatabase"
	metaKeyShardID   = "metaKeyShardID"
	metaKeyLeader    = "metaKeyLeader"
	metaKeyReplicas  = "metaKeyReplicas"
	metaKeyReplica   = "metaKeyReplica"
)

var (
	clientConnFct ClientConnFactory
)

func init() {
	clientConnFct = &clientConnFactory{
		connMap: make(map[models.Node]*grpc.ClientConn),
	}
}

// ClientConnFactory is the factory for grpc ClientConn.
type ClientConnFactory interface {
	// GetClientConn returns the grpc ClientConn for target node.
	// One connection for a target node.
	// Concurrent safe.
	GetClientConn(target models.Node) (*grpc.ClientConn, error)
}

// clientConnFactory implements ClientConnFactory.
type clientConnFactory struct {
	// target -> connection
	connMap map[models.Node]*grpc.ClientConn
	// lock to protect connMap
	lock4map sync.Mutex
}

// GetClientConnFactory returns a singleton ClientConnFactory.
func GetClientConnFactory() ClientConnFactory {
	return clientConnFct
}

// GetClientConn returns the grpc ClientConn for a target node.
// Concurrent safe.
func (fct *clientConnFactory) GetClientConn(target models.Node) (*grpc.ClientConn, error) {
	fct.lock4map.Lock()
	defer fct.lock4map.Unlock()

	coon, ok := fct.connMap[target]
	if ok {
		return coon, nil
	}
	conn, err := grpc.Dial(target.Indicator(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	fct.connMap[target] = conn

	return conn, nil
}

// ClientStreamFactory is the factory to get ClientStream.
type ClientStreamFactory interface {
	// LogicNode returns the a logic Node which will be transferred to the target server for identification.
	LogicNode() models.Node
	// CreateWriteClient creates a stream WriteClient.
	CreateWriteClient(db string, shardID int32, target models.Node) (storage.WriteService_WriteClient, error)
	// CreateQueryClient creates a stream task client
	CreateTaskClient(target models.Node) (common.TaskService_HandleClient, error)
	// CreateWriteServiceClient creates a WriteServiceClient
	CreateWriteServiceClient(target models.Node) (storage.WriteServiceClient, error)
	// CreateReplicaServiceClient creates a replicaRpc.ReplicaServiceClient.
	CreateReplicaServiceClient(target models.Node) (replicaRpc.ReplicaServiceClient, error)
}

// clientStreamFactory implements ClientStreamFactory.
type clientStreamFactory struct {
	logicNode models.Node
	connFct   ClientConnFactory
}

// NewClientStreamFactory returns a factory to get clientStream.
func NewClientStreamFactory(logicNode models.Node) ClientStreamFactory {
	return &clientStreamFactory{
		logicNode: logicNode,
		connFct:   GetClientConnFactory(),
	}
}

// LogicNode returns the a logic Node which will be transferred to the target server for identification.
func (w *clientStreamFactory) LogicNode() models.Node {
	return w.logicNode
}

// CreateQueryClient creates a stream task client
func (w *clientStreamFactory) CreateTaskClient(target models.Node) (common.TaskService_HandleClient, error) {
	conn, err := w.connFct.GetClientConn(target)
	if err != nil {
		return nil, err
	}

	node := w.LogicNode()
	//TODO handle context?????
	ctx := createOutgoingContextWithPairs(context.TODO(), metaKeyLogicNode, (&node).Indicator())
	cli, err := common.NewTaskServiceClient(conn).Handle(ctx)
	return cli, err
}

// CreateWriteClient creates a WriteClient.
func (w *clientStreamFactory) CreateWriteClient(db string, shardID int32,
	target models.Node) (storage.WriteService_WriteClient, error) {
	conn, err := w.connFct.GetClientConn(target)
	if err != nil {
		return nil, err
	}

	// pass logicNode.ID as meta to rpc serve
	ctx := createOutgoingContext(context.TODO(), db, shardID, w.LogicNode())
	cli, err := storage.NewWriteServiceClient(conn).Write(ctx)

	return cli, err
}

// CreateWriteServiceClient creates a WriteServiceClient
func (w *clientStreamFactory) CreateWriteServiceClient(target models.Node) (storage.WriteServiceClient, error) {
	conn, err := w.connFct.GetClientConn(target)
	if err != nil {
		return nil, err
	}
	return storage.NewWriteServiceClient(conn), nil
}

// CreateReplicaServiceClient creates a replicaRpc.ReplicaServiceClient.
func (w *clientStreamFactory) CreateReplicaServiceClient(target models.Node) (replicaRpc.ReplicaServiceClient, error) {
	conn, err := w.connFct.GetClientConn(target)
	if err != nil {
		return nil, err
	}
	return replicaRpc.NewReplicaServiceClient(conn), nil
}

// createOutgoingContextWithPairs creates outGoing context with key, value pairs.
func createOutgoingContextWithPairs(ctx context.Context, pairs ...string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(pairs...))
}

// createIncomingContextWithPairs creates outGoing context with key, value pairs, mainly for test.
func createIncomingContextWithPairs(ctx context.Context, pairs ...string) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.Pairs(pairs...))
}

// createOutgoingContext creates outGoing context with provided parameters.
// db is the database, shardID is the shard id for database,
// logicNode is a client provided identification on server side.
// These parameters will passed to the sever side in stream context.
func createOutgoingContext(ctx context.Context, db string, shardID int32, logicNode models.Node) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		metaKeyLogicNode, logicNode.Indicator(),
		metaKeyDatabase, db,
		metaKeyShardID, strconv.Itoa(int(shardID)))
}

// CreateIncomingContext creates incoming context with given parameters, mainly for test rpc server, mock incoming context.
func CreateIncomingContext(ctx context.Context, db string, shardID int32, logicNode models.Node) context.Context {
	return metadata.NewIncomingContext(ctx,
		metadata.Pairs(metaKeyLogicNode, logicNode.Indicator(),
			metaKeyDatabase, db,
			metaKeyShardID, strconv.Itoa(int(shardID))))
}

// CreateIncomingContextWithNode creates incoming context with given parameters, mainly for test rpc server, mock incoming context.
func CreateIncomingContextWithNode(ctx context.Context, node models.Node) context.Context {
	return createIncomingContextWithPairs(ctx, metaKeyLogicNode, node.Indicator())
}

// CreateOutgoingContextWithNode creates outgoing context with logic node.
func CreateOutgoingContextWithNode(ctx context.Context, node models.Node) context.Context {
	return createOutgoingContextWithPairs(ctx, metaKeyLogicNode, node.Indicator())
}

// getStringFromContext retrieving string metaValue from context for metaKey.
func getStringFromContext(ctx context.Context, metaKey string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("meta data not exists, key: " + metaKey)
	}

	strList := md.Get(metaKey)

	if len(strList) != 1 {
		return "", errors.New("meta data should have exactly one string value")
	}
	return strList[0], nil
}

// GetLogicNodeFromContext returns the logicNode.
func GetLogicNodeFromContext(ctx context.Context) (*models.Node, error) {
	strVal, err := getStringFromContext(ctx, metaKeyLogicNode)
	if err != nil {
		return nil, err
	}

	return models.ParseNode(strVal)
}

// GetDatabaseFromContext returns database.
func GetDatabaseFromContext(ctx context.Context) (string, error) {
	return getStringFromContext(ctx, metaKeyDatabase)
}

// GetShardIDFromContext returns shardID.
func GetShardIDFromContext(ctx context.Context) (int32, error) {
	return getIntFromContext(ctx, metaKeyShardID)
}

// GetLeaderFromContext returns leader's node id.
func GetLeaderFromContext(ctx context.Context) (models.NodeID, error) {
	nodeID, err := getIntFromContext(ctx, metaKeyLeader)
	if err != nil {
		return models.NodeID(-1), err
	}
	return models.NodeID(nodeID), nil
}

// GetFollowerFromContext returns follower's node id.
func GetFollowerFromContext(ctx context.Context) (models.NodeID, error) {
	nodeID, err := getIntFromContext(ctx, metaKeyReplica)
	if err != nil {
		return models.NodeID(-1), err
	}
	return models.NodeID(nodeID), nil
}

// GetReplicasFromContext returns replicas' node id.
func GetReplicasFromContext(ctx context.Context) ([]models.NodeID, error) {
	nodeIDs, err := getStringFromContext(ctx, metaKeyReplicas)
	if err != nil {
		return nil, err
	}
	var replicas []models.NodeID
	if err := json.Unmarshal([]byte(nodeIDs), &replicas); err != nil {
		return nil, err
	}
	return replicas, nil
}

// getIntFromContext retrieving int metaValue from context for metaKey.
func getIntFromContext(ctx context.Context, metaKey string) (int32, error) {
	strVal, err := getStringFromContext(ctx, metaKey)
	if err != nil {
		return -1, err
	}

	num, err := strconv.Atoi(strVal)
	if err != nil {
		return -1, err
	}

	return int32(num), nil
}
