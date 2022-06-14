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
	"errors"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/conntrack"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
)

//go:generate mockgen -source ./rpc.go -destination=./rpc_mock.go -package=rpc

var (
	brokerClientConnFct  ClientConnFactory
	storageClientConnFct ClientConnFactory
)

func init() {
	brokerClientConnFct = &clientConnFactory{
		connMap:       make(map[string]*grpc.ClientConn),
		clientTracker: conntrack.NewGRPCClientTracker(linmetric.BrokerRegistry),
	}
	storageClientConnFct = &clientConnFactory{
		connMap:       make(map[string]*grpc.ClientConn),
		clientTracker: conntrack.NewGRPCClientTracker(linmetric.StorageRegistry),
	}
}

// ClientConnFactory is the factory for grpc ClientConn.
type ClientConnFactory interface {
	// GetClientConn returns the grpc ClientConn for target node.
	// One connection for a target node.
	// Concurrent safe.
	GetClientConn(target models.Node) (*grpc.ClientConn, error)
	// CloseClientConn closes client connection for spec target node.
	CloseClientConn(target models.Node) error
}

// clientConnFactory implements ClientConnFactory.
type clientConnFactory struct {
	// target's indicator -> connection
	connMap map[string]*grpc.ClientConn
	// lock to protect connMap
	mu            sync.RWMutex
	clientTracker *conntrack.GRPCClientTracker
}

// GetBrokerClientConnFactory returns a singleton ClientConnFactory for broker side.
func GetBrokerClientConnFactory() ClientConnFactory {
	return brokerClientConnFct
}

// GetStorageClientConnFactory returns a singleton ClientConnFactory for storage side.
func GetStorageClientConnFactory() ClientConnFactory {
	return storageClientConnFct
}

// GetClientConn returns the grpc ClientConn for a target node.
// Concurrent safe.
func (fct *clientConnFactory) GetClientConn(target models.Node) (*grpc.ClientConn, error) {
	indicator := target.Indicator()
	fct.mu.RLock()
	conn, ok := fct.connMap[indicator]
	fct.mu.RUnlock()
	if ok {
		return conn, nil
	}

	fct.mu.Lock()
	defer fct.mu.Unlock()

	// double check
	conn, ok = fct.connMap[indicator]
	if ok {
		return conn, nil
	}
	conn, err := grpc.Dial(
		target.Indicator(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(fct.clientTracker.StreamClientInterceptor()),
		grpc.WithUnaryInterceptor(fct.clientTracker.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	fct.connMap[indicator] = conn
	return conn, nil
}

// CloseClientConn closes client connection for spec target node.
func (fct *clientConnFactory) CloseClientConn(target models.Node) error {
	indicator := target.Indicator()

	fct.mu.RLock()
	conn, ok := fct.connMap[indicator]
	fct.mu.RUnlock()

	if ok {
		// if connection exist for node
		if err := conn.Close(); err != nil {
			// if close err, keep it, try reconnect, maybe get some err for connection closed before reconnected
			return err
		}
		// if close success, need remove connection from cache
		fct.mu.Lock()
		delete(fct.connMap, indicator)
		fct.mu.Unlock()
	}

	return nil
}

// ClientStreamFactory is the factory to get ClientStream.
type ClientStreamFactory interface {
	// LogicNode returns the a logic Node which will be transferred to the target server for identification.
	LogicNode() models.Node
	// CreateTaskClient creates a stream task client
	CreateTaskClient(target models.Node) (protoCommonV1.TaskService_HandleClient, error)
	// CreateReplicaServiceClient creates a protoReplicaV1.ReplicaServiceClient.
	CreateReplicaServiceClient(target models.Node) (protoReplicaV1.ReplicaServiceClient, error)
	// CreateWriteServiceClient creates a protoWriteV1.WriteServiceClient.
	CreateWriteServiceClient(target models.Node) (protoWriteV1.WriteServiceClient, error)
}

// clientStreamFactory implements ClientStreamFactory.
type clientStreamFactory struct {
	ctx       context.Context
	logicNode models.Node
	connFct   ClientConnFactory
}

// NewClientStreamFactory returns a factory to get clientStream.
func NewClientStreamFactory(ctx context.Context, logicNode models.Node, connFct ClientConnFactory) ClientStreamFactory {
	return &clientStreamFactory{
		ctx:       ctx,
		logicNode: logicNode,
		connFct:   connFct,
	}
}

// LogicNode returns the a logic Node which will be transferred to the target server for identification.
func (w *clientStreamFactory) LogicNode() models.Node {
	return w.logicNode
}

// CreateTaskClient creates a stream task client
func (w *clientStreamFactory) CreateTaskClient(target models.Node) (protoCommonV1.TaskService_HandleClient, error) {
	conn, err := w.connFct.GetClientConn(target)
	if err != nil {
		return nil, err
	}

	node := w.LogicNode()
	// https://pkg.go.dev/google.golang.org/grpc#ClientConn.NewStream
	// context is the lifetime of stream
	ctx := CreateOutgoingContextWithPairs(w.ctx, constants.RPCMetaKeyLogicNode, node.Indicator())
	cli, err := protoCommonV1.NewTaskServiceClient(conn).Handle(ctx)
	return cli, err
}

// CreateReplicaServiceClient creates a protoReplicaV1.ReplicaServiceClient.
func (w *clientStreamFactory) CreateReplicaServiceClient(target models.Node) (protoReplicaV1.ReplicaServiceClient, error) {
	conn, err := w.connFct.GetClientConn(target)
	if err != nil {
		return nil, err
	}
	return protoReplicaV1.NewReplicaServiceClient(conn), nil
}

// CreateWriteServiceClient creates a protoWriteV1.WriteServiceClient.
func (w *clientStreamFactory) CreateWriteServiceClient(target models.Node) (protoWriteV1.WriteServiceClient, error) {
	conn, err := w.connFct.GetClientConn(target)
	if err != nil {
		return nil, err
	}
	return protoWriteV1.NewWriteServiceClient(conn), nil
}

// CreateOutgoingContextWithPairs creates outGoing context with key, value pairs.
func CreateOutgoingContextWithPairs(ctx context.Context, pairs ...string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(pairs...))
}

// GetStringFromContext retrieving string metaValue from context for metaKey.
func GetStringFromContext(ctx context.Context, metaKey string) (string, error) {
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
func GetLogicNodeFromContext(ctx context.Context) (models.Node, error) {
	strVal, err := GetStringFromContext(ctx, constants.RPCMetaKeyLogicNode)
	if err != nil {
		return nil, err
	}

	return models.ParseNode(strVal)
}
