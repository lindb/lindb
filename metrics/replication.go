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

package metrics

import "github.com/lindb/lindb/internal/linmetric"

// BrokerDatabaseWriteStatistics represents database channel write statistics.
type BrokerDatabaseWriteStatistics struct {
	OutOfTimeRange *linmetric.BoundCounter // timestamp of metrics out of acceptable write time range
	ShardNotFound  *linmetric.BoundCounter // shard not found count
}

// BrokerFamilyWriteStatistics represents family channel write statistics.
type BrokerFamilyWriteStatistics struct {
	ActiveWriteFamilies  *linmetric.BoundGauge   // number of current active replica family channel
	BatchMetrics         *linmetric.BoundCounter // batch into memory chunk success count
	BatchMetricFailures  *linmetric.BoundCounter // batch into memory chunk failure count
	PendingSend          *linmetric.BoundGauge   // number of pending send message
	SendSuccess          *linmetric.BoundCounter // send message success count
	SendFailure          *linmetric.BoundCounter // send message failure count
	SendSize             *linmetric.BoundCounter // bytes of send message
	Retry                *linmetric.BoundCounter // retry count
	RetryDrop            *linmetric.BoundCounter // number of drop message after too many retry
	CreateStream         *linmetric.BoundCounter // create replica stream success count
	CreateStreamFailures *linmetric.BoundCounter // create replica stream failure count
	CloseStream          *linmetric.BoundCounter // close replica stream success count
	CloseStreamFailures  *linmetric.BoundCounter // close replica stream failure count
	LeaderChanged        *linmetric.BoundCounter // shard leader changed
}

// StorageLocalReplicatorStatistics represents local replicator statistics.
type StorageLocalReplicatorStatistics struct {
	DecompressFailures *linmetric.BoundCounter // decompress message failure count
	ReplicaFailures    *linmetric.BoundCounter // replica failure count
	ReplicaRows        *linmetric.BoundCounter // row number of replica
	AckSequence        *linmetric.BoundCounter // ack persist sequence count
	InvalidSequence    *linmetric.BoundCounter // invalid replica sequence count
}

// StorageRemoteReplicatorStatistics represents remote replicator statistics.
type StorageRemoteReplicatorStatistics struct {
	NotReady                       *linmetric.BoundCounter // remote replicator channel not ready
	FollowerOffline                *linmetric.BoundCounter // remote follower node offline
	NeedCloseLastStream            *linmetric.BoundCounter // need close last stream, when do re-connection
	CloseLastStreamFailures        *linmetric.BoundCounter // close last stream failure
	CreateReplicaCli               *linmetric.BoundCounter // create replica client success
	CreateReplicaCliFailures       *linmetric.BoundCounter // create replica client failure
	CreateReplicaStream            *linmetric.BoundCounter // create replica stream success
	CreateReplicaStreamFailures    *linmetric.BoundCounter // create replica stream failure
	GetLastAckFailures             *linmetric.BoundCounter // get last ack sequence from remote follower failure
	ResetFollowerAppendIdx         *linmetric.BoundCounter // reset follower append index success
	ResetFollowerAppendIdxFailures *linmetric.BoundCounter // reset follower append index failure
	ResetAppendIdx                 *linmetric.BoundCounter // reset current leader local append index
	ResetReplicaIdx                *linmetric.BoundCounter // reset current leader replica index success
	ResetReplicaIdxFailures        *linmetric.BoundCounter // reset current leader replica index failure
	SendMsg                        *linmetric.BoundCounter // send replica msg success
	SendMsgFailures                *linmetric.BoundCounter // send replica msg failure
	ReceiveMsg                     *linmetric.BoundCounter // receive replica resp success
	ReceiveMsgFailures             *linmetric.BoundCounter // receive replica resp failure
	AckSequence                    *linmetric.BoundCounter // ack replica successfully sequence count
	InvalidAckSequence             *linmetric.BoundCounter // get wrong replica ack sequence from follower
}

// StorageReplicatorRunnerStatistics represents storage replicator runner statistics.
type StorageReplicatorRunnerStatistics struct {
	ActiveReplicators      *linmetric.BoundGauge   // number of current active local replicator
	ReplicaPanics          *linmetric.BoundCounter // replica panic count
	ConsumeMessage         *linmetric.BoundCounter // get message success count
	ConsumeMessageFailures *linmetric.BoundCounter // get message failure count
	ReplicaLag             *linmetric.BoundGauge   // replica lag message count
	ReplicaBytes           *linmetric.BoundCounter // bytes of replica data
	Replica                *linmetric.BoundCounter // replica success count
}

// StorageWriteAheadLogStatistics represents storage write ahead log statistics.
type StorageWriteAheadLogStatistics struct {
	ReceiveWriteSize   *linmetric.BoundCounter // receive write request bytes(broker->leader)
	WriteWAL           *linmetric.BoundCounter // write wal success(broker->leader)
	WriteWALFailures   *linmetric.BoundCounter // write wal failure(broker->leader)
	ReceiveReplicaSize *linmetric.BoundCounter // receive replica request bytes(storage leader->follower)
	ReplicaWAL         *linmetric.BoundCounter // replica wal success(storage leader->follower)
	ReplicaWALFailures *linmetric.BoundCounter // replica wal failure(storage leader->follower)
}

// NewBrokerDatabaseWriteStatistics creates a database channel write statistics.
func NewBrokerDatabaseWriteStatistics(database string) *BrokerDatabaseWriteStatistics {
	scope := linmetric.BrokerRegistry.NewScope("lindb.broker.database.write")
	return &BrokerDatabaseWriteStatistics{
		OutOfTimeRange: scope.NewCounterVec("out_of_time_range", "db").WithTagValues(database),
		ShardNotFound:  scope.NewCounterVec("shard_not_found", "db").WithTagValues(database),
	}
}

// NewBrokerFamilyWriteStatistics creates a family channel write statistics.
func NewBrokerFamilyWriteStatistics(database string) *BrokerFamilyWriteStatistics {
	scope := linmetric.BrokerRegistry.NewScope("lindb.broker.family.write")
	return &BrokerFamilyWriteStatistics{
		ActiveWriteFamilies:  scope.NewGaugeVec("active_families", "db").WithTagValues(database),
		BatchMetrics:         scope.NewCounterVec("batch_metrics", "db").WithTagValues(database),
		BatchMetricFailures:  scope.NewCounterVec("batch_metrics_failures", "db").WithTagValues(database),
		PendingSend:          scope.NewGaugeVec("pending_send", "db").WithTagValues(database),
		SendSuccess:          scope.NewCounterVec("send_success", "db").WithTagValues(database),
		SendFailure:          scope.NewCounterVec("send_failures", "db").WithTagValues(database),
		SendSize:             scope.NewCounterVec("send_size", "db").WithTagValues(database),
		Retry:                scope.NewCounterVec("retry", "db").WithTagValues(database),
		RetryDrop:            scope.NewCounterVec("retry_drop", "db").WithTagValues(database),
		CreateStream:         scope.NewCounterVec("create_stream", "db").WithTagValues(database),
		CreateStreamFailures: scope.NewCounterVec("create_stream_failures", "db").WithTagValues(database),
		CloseStream:          scope.NewCounterVec("close_stream", "db").WithTagValues(database),
		CloseStreamFailures:  scope.NewCounterVec("close_stream_failures", "db").WithTagValues(database),
		LeaderChanged:        scope.NewCounterVec("leader_changed", "db").WithTagValues(database),
	}
}

// NewStorageLocalReplicatorStatistics creates a storage local replicator statistics.
func NewStorageLocalReplicatorStatistics(database, shard string) *StorageLocalReplicatorStatistics {
	scope := linmetric.StorageRegistry.NewScope("lindb.storage.replica.local")
	return &StorageLocalReplicatorStatistics{
		DecompressFailures: scope.NewCounterVec("decompress_failures", "db", "shard").WithTagValues(database, shard),
		ReplicaFailures:    scope.NewCounterVec("replica_failures", "db", "shard").WithTagValues(database, shard),
		ReplicaRows:        scope.NewCounterVec("replica_rows", "db", "shard").WithTagValues(database, shard),
		AckSequence:        scope.NewCounterVec("ack_sequence", "db", "shard").WithTagValues(database, shard),
		InvalidSequence:    scope.NewCounterVec("invalid_sequence", "db", "shard").WithTagValues(database, shard),
	}
}

// NewStorageReplicatorRunnerStatistics creates storage replicator runner statistics.
func NewStorageReplicatorRunnerStatistics(replicatorType, database, shard string) *StorageReplicatorRunnerStatistics {
	scope := linmetric.StorageRegistry.NewScope("lindb.storage.replicator.runner")
	return &StorageReplicatorRunnerStatistics{
		ActiveReplicators: scope.NewGaugeVec("active_replicators", "type", "db", "shard").
			WithTagValues(replicatorType, database, shard),
		ReplicaPanics: scope.NewCounterVec("replica_panics", "type", "db", "shard").
			WithTagValues(replicatorType, database, shard),
		ConsumeMessage: scope.NewCounterVec("consume_msg", "type", "db", "shard").
			WithTagValues(replicatorType, database, shard),
		ConsumeMessageFailures: scope.NewCounterVec("consume_msg_failures", "type", "db", "shard").
			WithTagValues(replicatorType, database, shard),
		ReplicaLag: scope.NewGaugeVec("replica_lag", "type", "db", "shard").
			WithTagValues(replicatorType, database, shard),
		ReplicaBytes: scope.NewCounterVec("replica_bytes", "type", "db", "shard").
			WithTagValues(replicatorType, database, shard),
		Replica: scope.NewCounterVec("replicas", "type", "db", "shard").
			WithTagValues(replicatorType, database, shard),
	}
}

// NewStorageRemoteReplicatorStatistics creates remote replicator statistics.
func NewStorageRemoteReplicatorStatistics(database, shard string) *StorageRemoteReplicatorStatistics {
	scope := linmetric.StorageRegistry.NewScope("lindb.storage.replica.remote")
	return &StorageRemoteReplicatorStatistics{
		NotReady: scope.NewCounterVec("not_ready", "db", "shard").
			WithTagValues(database, shard),
		FollowerOffline: scope.NewCounterVec("follower_offline", "db", "shard").
			WithTagValues(database, shard),
		NeedCloseLastStream: scope.NewCounterVec("need_close_last_stream", "db", "shard").
			WithTagValues(database, shard),
		CloseLastStreamFailures: scope.NewCounterVec("close_last_stream_failures", "db", "shard").
			WithTagValues(database, shard),
		CreateReplicaCli: scope.NewCounterVec("create_replica_cli", "db", "shard").
			WithTagValues(database, shard),
		CreateReplicaCliFailures: scope.NewCounterVec("create_replica_cli_failures", "db", "shard").
			WithTagValues(database, shard),
		CreateReplicaStream: scope.NewCounterVec("create_replica_stream", "db", "shard").
			WithTagValues(database, shard),
		CreateReplicaStreamFailures: scope.NewCounterVec("create_replica_stream_failures", "db", "shard").
			WithTagValues(database, shard),
		GetLastAckFailures: scope.NewCounterVec("get_last_ack_failures", "db", "shard").
			WithTagValues(database, shard),
		ResetFollowerAppendIdx: scope.NewCounterVec("reset_follower_append_idx", "db", "shard").
			WithTagValues(database, shard),
		ResetFollowerAppendIdxFailures: scope.NewCounterVec("reset_follower_append_idx_failures", "db", "shard").
			WithTagValues(database, shard),
		ResetAppendIdx: scope.NewCounterVec("reset_append_idx", "db", "shard").
			WithTagValues(database, shard),
		ResetReplicaIdx: scope.NewCounterVec("reset_replica_idx", "db", "shard").
			WithTagValues(database, shard),
		ResetReplicaIdxFailures: scope.NewCounterVec("reset_replica_failures", "db", "shard").
			WithTagValues(database, shard),
		SendMsg: scope.NewCounterVec("send_msg", "db", "shard").
			WithTagValues(database, shard),
		SendMsgFailures: scope.NewCounterVec("send_msg_failures", "db", "shard").
			WithTagValues(database, shard),
		ReceiveMsg: scope.NewCounterVec("receive_msg", "db", "shard").
			WithTagValues(database, shard),
		ReceiveMsgFailures: scope.NewCounterVec("receive_msg_failures", "db", "shard").
			WithTagValues(database, shard),
		AckSequence: scope.NewCounterVec("ack_sequence", "db", "shard").
			WithTagValues(database, shard),
		InvalidAckSequence: scope.NewCounterVec("invalid_ack_sequence", "db", "shard").
			WithTagValues(database, shard),
	}
}

// NewStorageWriteAheadLogStatistics creates a storage write ahead log statistics.
func NewStorageWriteAheadLogStatistics(database, shard string) *StorageWriteAheadLogStatistics {
	scope := linmetric.StorageRegistry.NewScope("lindb.storage.wal")
	return &StorageWriteAheadLogStatistics{
		ReceiveWriteSize: scope.NewCounterVec("receive_write_bytes", "db", "shard").
			WithTagValues(database, shard),
		WriteWAL: scope.NewCounterVec("write_wal", "db", "shard").
			WithTagValues(database, shard),
		WriteWALFailures: scope.NewCounterVec("write_wal_failures", "db", "shard").
			WithTagValues(database, shard),
		ReceiveReplicaSize: scope.NewCounterVec("receive_replica_bytes", "db", "shard").
			WithTagValues(database, shard),
		ReplicaWAL: scope.NewCounterVec("replica_wal", "db", "shard").
			WithTagValues(database, shard),
		ReplicaWALFailures: scope.NewCounterVec("replica_wal_failures", "db", "shard").
			WithTagValues(database, shard),
	}
}
