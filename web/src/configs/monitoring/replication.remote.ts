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
import { MonitoringDB } from "@src/constants";
import { Dashboard, Unit } from "@src/models";
import { chartOptions } from "./system";

export const RemoteReplicationDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Replica Lag",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'replica_lag' from 'lindb.storage.replicator.runner' where type='remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 24,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Number Of Replica",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'replicas' from 'lindb.storage.replicator.runner' where type='remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Replica Traffic",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'replica_bytes' from 'lindb.storage.replicator.runner' where type='remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Bytes,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Active Replicators",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'active_replicators' from 'lindb.storage.replicator.runner' where type='remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Consumer Message",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'consume_msg' from 'lindb.storage.replicator.runner' where type='remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Consumer Message Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'consume_msg_failures' from 'lindb.storage.replicator.runner' where type='remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Replica Painc",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'replica_panics' from 'lindb.storage.replicator.runner' where type='remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Send Message",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'send_msg' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Send Message Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'send_msg_failures' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Receive Message",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'receive_msg' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Receive Message Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'receive_msg_failures' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Ack Replica Sequence",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'ack_sequence' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Invalid Ack Sequence",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'invalid_ack_sequence' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Not Ready",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'not_ready' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Remote Follower Offline",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'follower_offline' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Need Close Last Replica Stream",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'need_close_last_stream' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Close Last Stream Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'close_last_stream_failures' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Create Replica Client",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'create_replica_cli' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Create Reaplica Client Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'create_replica_cli_failures' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Create Replica Stream",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'create_replica_stream' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Create Replica Stream Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'create_replica_stream_failures' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Get Last Ack Sequence Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'get_last_ack_failures' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Reset Remote Follower Append Index",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'reset_follower_append_idx' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Reset Remote Follower Append Index Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'reset_follower_append_idx_failures' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Reset Append Index",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'reset_append_idx' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Reset Replica Index",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'reset_replica_idx' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Reset Replica Index Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'reset_replica_failures' from 'lindb.storage.replica.remote' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
      ],
    },
  ],
};
