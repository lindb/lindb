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

package query

import (
	"fmt"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
)

// transportManager implments rpc.TransportManager interface.
type transportManager struct {
	taskClientFactory rpc.TaskClientFactory
	taskServerFactory rpc.TaskServerFactory

	statistics *metrics.TransportStatistics

	logger logger.Logger
}

// NewTransportManager creates a rpc transport manager instance.
func NewTransportManager(
	taskClientFactory rpc.TaskClientFactory,
	taskServerFactory rpc.TaskServerFactory,
	registry *linmetric.Registry,
) rpc.TransportManager {
	return &transportManager{
		taskClientFactory: taskClientFactory,
		taskServerFactory: taskServerFactory,
		statistics:        metrics.NewTransportStatistics(registry),
		logger:            logger.GetLogger("Query", "TransportManager"),
	}
}

// SendRequest sends the task request to target node.
func (mgr *transportManager) SendRequest(targetNodeID string, req *protoCommonV1.TaskRequest) error {
	mgr.logger.Debug("send query task", logger.String("target", targetNodeID))
	client := mgr.taskClientFactory.GetTaskClient(targetNodeID)
	if client == nil {
		mgr.statistics.SentRequestFailures.Incr()
		return fmt.Errorf("SendRequest: %w, targetNodeID: %s", ErrNoSendStream, targetNodeID)
	}
	if err := client.Send(req); err != nil {
		mgr.statistics.SentRequestFailures.Incr()
		return fmt.Errorf("SendRequest: %w, targetNodeID: %s", ErrTaskSend, targetNodeID)
	}
	mgr.statistics.SentRequest.Incr()
	return nil
}

// SendResponse sends the task response to target node.
func (mgr *transportManager) SendResponse(targetNodeID string, resp *protoCommonV1.TaskResponse) error {
	stream := mgr.taskServerFactory.GetStream(targetNodeID)
	if stream == nil {
		mgr.statistics.SentResponseFailures.Incr()
		return fmt.Errorf("SendResponse: %w, parentNodeID: %s", ErrNoSendStream, targetNodeID)
	}
	if err := stream.Send(resp); err != nil {
		mgr.statistics.SentResponseFailures.Incr()
		return fmt.Errorf("SendResponse: %w, parentNodeID: %s", ErrResponseSend, targetNodeID)
	}
	mgr.statistics.SentResponses.Incr()
	return nil
}
