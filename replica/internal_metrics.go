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

import "github.com/lindb/lindb/internal/linmetric"

var (
	walScope             = linmetric.StorageRegistry.NewScope("lindb.wal")
	appendSeqVec         = walScope.Scope("write").NewGaugeVec("append_seq", "db", "shard", "family")
	replicaSeqVec        = walScope.Scope("replica").NewGaugeVec("consume_seq", "db", "shard", "family", "from", "to")
	activeReplicaChannel = walScope.NewGaugeVec("active_replicas", "db", "type")
	receiveWriteSize     = walScope.NewCounterVec("receive_write_size", "db")
	writeWAL             = walScope.NewCounterVec("write_wal", "db")
	writeWALFailure      = walScope.NewCounterVec("write_wal_failure", "db")
	receiveReplicaSize   = walScope.NewCounterVec("receive_replica_size", "db")
	replicaWAL           = walScope.NewCounterVec("replica_wal", "db")
	replicaWALFailure    = walScope.NewCounterVec("replica_wal_failure", "db")

	localReplicaScope       = linmetric.StorageRegistry.NewScope("lindb.replica.local")
	localMaxDecodedBlockVec = localReplicaScope.NewMaxVec("max_decoded_block", "db", "shard")
	localReplicaCountsVec   = localReplicaScope.NewCounterVec("replica_count", "db", "shard")
	localReplicaBytesVec    = localReplicaScope.NewCounterVec("replica_bytes", "db", "shard")
	localReplicaRowsVec     = localReplicaScope.NewCounterVec("replica_rows", "db", "shard")
	localReplicaSequenceVec = localReplicaScope.NewGaugeVec("replica_sequence", "db", "shard")
	localInvalidSequenceVec = localReplicaScope.NewCounterVec("invalid_sequence", "db", "shard")
)

var (
	brokerScope         = linmetric.BrokerRegistry.NewScope("lindb.broker.replica")
	activeWriteFamilies = brokerScope.NewGaugeVec("active_families", "db")
	batchMetrics        = brokerScope.NewCounterVec("batch_metrics", "db")
	batchMetricFailures = brokerScope.NewCounterVec("batch_metrics_failures", "db")
	pendingSend         = brokerScope.NewGaugeVec("pending_send", "db")
	sendSuccess         = brokerScope.NewCounterVec("send_success", "db")
	sendFailure         = brokerScope.NewCounterVec("send_failure", "db")
	sendSize            = brokerScope.NewCounterVec("send_size", "db")
	retryCount          = brokerScope.NewCounterVec("retry", "db")
	retryDrop           = brokerScope.NewCounterVec("retry_drop", "db")
)
