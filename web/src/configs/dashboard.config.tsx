/*
Licensed to LinDB under one or more contributor
license agreements. See the NOTICE file distributed with
this work for additional information regarding copyright
ownership. LinDB licenses this file to you under
the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
import {
  RuntimeDashboard,
  SystemDashboard,
  ConcurrentPoolDashboard,
  ConcurrentLimitDashboard,
  IngestionDashboard,
  BrokerWriteDashboard,
  BrokerCoordinatorDashboard,
  MasterCoordinatorDashboard,
  MasterControllerDashboard,
  StorageCoordinatorDashboard,
  BrokerQueryDashboard,
  StorageQueryDashboard,
  KVStoreReadDashboard,
  KVStoreWriteDashboard,
  KVStoreJobDashboard,
  NetworkTCPDasbhoard,
  NetworkGRPCDasbhoard,
  TSDBWriteDashboard,
  TSDBJobDashboard,
  WALDashboard,
  LocalReplicationDashboard,
  RemoteReplicationDashboard,
} from "@src/configs";
import { MonitoringDB, StateRoleName } from "@src/constants";
import { DashboardItem, Variate } from "@src/models";

export enum DashboardName {
  Runtime = "Runtime",
  System = "System",
  KVRead = "KV Read",
  KVWrite = "KV Write",
  KVJob = "KV Job",
  TSDBWrite = "TSDB Write",
  TSDBJob = "TSDB Job",
  ConcurrentPool = "Concurrent Pool",
  ConcurrentLimit = "Concurrent Limit",
  CoordinatorMaster = "Master Coordinator",
  MasterController = "Master Controller",
  CoordinatorBroker = "Broker Coordinator",
  CoordinatorStorage = "Storage Coordinator",
  Ingestion = "Ingestion",
  BrokerWrite = "Broker Write",
  BrokerQuery = "Broker Query",
  StorageQuery = "Storage Query",
  NetworkTCP = "Network TCP",
  NetworkGRPC = "Network GRPC",
  WAL = "WAL",
  LocalReplication = "Local Replication",
  RemoteReplication = "Remote Replication",
}

export const CommonVariates = [
  {
    tagKey: "namespace",
    label: "Namespace",
    watch: { cascade: ["role"], clear: "node" },
    scope: [StateRoleName.Storage],
    db: MonitoringDB,
    sql: "show tag values from 'lindb.runtime.mem' with key=namespace",
  },
  {
    tagKey: "node",
    label: "Node",
    watch: { cascade: ["namespace", "role"] },
    db: MonitoringDB,
    scope: [StateRoleName.Storage, StateRoleName.Broker],
    multiple: true,
    sql: "show tag values from 'lindb.runtime.mem' with key=node",
  },
] as Variate[];

export const Dashboards = [
  {
    value: "broker.query",
    label: DashboardName.BrokerQuery,
    scope: [StateRoleName.Broker],
    dashboard: BrokerQueryDashboard,
  },
  {
    value: "ingestion",
    label: DashboardName.Ingestion,
    scope: [StateRoleName.Broker],
    dashboard: IngestionDashboard,
  },
  {
    value: "broker.write",
    label: DashboardName.BrokerWrite,
    scope: [StateRoleName.Broker],
    dashboard: BrokerWriteDashboard,
  },
  {
    value: "master.controller",
    label: DashboardName.MasterController,
    scope: [StateRoleName.Broker],
    dashboard: MasterControllerDashboard,
  },
  {
    value: "coordinator.master",
    label: DashboardName.CoordinatorMaster,
    scope: [StateRoleName.Broker],
    dashboard: MasterCoordinatorDashboard,
  },
  {
    value: "coordinator.broker",
    label: DashboardName.CoordinatorBroker,
    scope: [StateRoleName.Broker],
    dashboard: BrokerCoordinatorDashboard,
  },
  {
    value: "storage.query",
    label: DashboardName.StorageQuery,
    scope: [StateRoleName.Storage],
    dashboard: StorageQueryDashboard,
  },
  {
    value: "wal",
    label: DashboardName.WAL,
    scope: [StateRoleName.Storage],
    dashboard: WALDashboard,
  },
  {
    value: "local.replication",
    label: DashboardName.LocalReplication,
    scope: [StateRoleName.Storage],
    dashboard: LocalReplicationDashboard,
  },
  {
    value: "remote.replication",
    label: DashboardName.RemoteReplication,
    scope: [StateRoleName.Storage],
    dashboard: RemoteReplicationDashboard,
  },
  {
    value: "coordinator.storage",
    label: DashboardName.CoordinatorStorage,
    scope: [StateRoleName.Storage],
    dashboard: StorageCoordinatorDashboard,
  },
  {
    value: "kv.read",
    label: DashboardName.KVRead,
    scope: [StateRoleName.Storage],
    dashboard: KVStoreReadDashboard,
  },
  {
    value: "kv.write",
    label: DashboardName.KVWrite,
    scope: [StateRoleName.Storage],
    dashboard: KVStoreWriteDashboard,
  },
  {
    value: "kv.job",
    label: DashboardName.KVJob,
    scope: [StateRoleName.Storage],
    dashboard: KVStoreJobDashboard,
  },
  {
    value: "tsdb.write",
    label: DashboardName.TSDBWrite,
    scope: [StateRoleName.Storage],
    dashboard: TSDBWriteDashboard,
  },
  {
    value: "tsdb.job",
    label: DashboardName.TSDBJob,
    scope: [StateRoleName.Storage],
    dashboard: TSDBJobDashboard,
  },
  {
    value: "concurrnet.pool",
    label: DashboardName.ConcurrentPool,
    scope: [StateRoleName.Storage, StateRoleName.Broker],
    dashboard: ConcurrentPoolDashboard,
  },
  {
    value: "concurrnet.limit",
    label: DashboardName.ConcurrentLimit,
    scope: [StateRoleName.Broker],
    dashboard: ConcurrentLimitDashboard,
  },
  {
    value: "network.tcp",
    label: DashboardName.NetworkTCP,
    scope: [StateRoleName.Storage, StateRoleName.Broker],
    dashboard: NetworkTCPDasbhoard,
  },
  {
    value: "newwork.grpc",
    label: DashboardName.NetworkGRPC,
    scope: [StateRoleName.Storage, StateRoleName.Broker],
    dashboard: NetworkGRPCDasbhoard,
  },
  {
    value: "runtime",
    label: DashboardName.Runtime,
    scope: [StateRoleName.Storage, StateRoleName.Broker],
    dashboard: RuntimeDashboard,
  },
  {
    value: "system",
    label: DashboardName.System,
    scope: [StateRoleName.Storage, StateRoleName.Broker],
    dashboard: SystemDashboard,
  },
] as DashboardItem[];
