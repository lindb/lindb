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

package handler

import (
	"context"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/service"

	replicaRpc "github.com/lindb/lindb/rpc/proto/replica"
)

// replicaHandler implements replica.ReplicaServiceServer interface for handling replica rpc request.
type replicaHandler struct {
	storageService service.StorageService

	logger *logger.Logger
}

// NewReplicaHandler creates a replica handler.
func NewReplicaHandler(storageService service.StorageService) replicaRpc.ReplicaServiceServer {
	return &replicaHandler{
		storageService: storageService,
		logger:         logger.GetLogger("storage", "replicaRpc"),
	}
}

// GetReplicaAckIndex returns current replica ack index.
func (r *replicaHandler) GetReplicaAckIndex(ctx context.Context,
	req *replicaRpc.GetReplicaAckIndexRequest,
) (*replicaRpc.GetReplicaAckIndexResponse, error) {
	panic("implement me")
}

// Reset resets replica index.
func (r *replicaHandler) Reset(ctx context.Context,
	request *replicaRpc.ResetIndexRequest,
) (*replicaRpc.ResetIndexResponse, error) {
	panic("implement me")
}

// Replica does replica request, and writes data.
func (r *replicaHandler) Replica(server replicaRpc.ReplicaService_ReplicaServer) error {
	panic("implement me")
}
