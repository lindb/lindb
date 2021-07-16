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

import pb "github.com/lindb/lindb/proto/gen/v1/common"

//go:generate mockgen -source ./task_receiver.go -destination=./task_receiver_mock.go -package=rpc

// TaskReceiver represents the task result receiver
type TaskReceiver interface {
	// Receive receives the task result
	Receive(req *pb.TaskResponse) error
}
