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

package replica

import "errors"

var (
	// define error types
	errChannelNotFound = errors.New("shard replica channel not found")
	errInvalidShardID  = errors.New("numOfShard should be greater than 0 and shardID should less then numOfShard")
	errInvalidShardNum = errors.New("numOfShard should be equal or greater than original setting")
	// ErrFamilyChannelCanceled is the error returned when a family channel is closed.
	ErrFamilyChannelCanceled = errors.New("family Channel is canceled")
	ErrIngestTimeout         = errors.New("ingest timout")
)
