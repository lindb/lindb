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

export const LocalReplicationDashboard: Dashboard = {
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
                sql: "select 'replica_lag' from 'lindb.storage.replicator.runner' where type='local' group by db,node",
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
                sql: "select 'replicas' from 'lindb.storage.replicator.runner' where type='local' group by db,node",
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
                sql: "select 'replica_bytes' from 'lindb.storage.replicator.runner' where type='local' group by db,node",
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
                sql: "select 'active_replicators' from 'lindb.storage.replicator.runner' where type='local' group by db,node",
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
                sql: "select 'consume_msg' from 'lindb.storage.replicator.runner' where type='local' group by db,node",
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
                sql: "select 'consume_msg_failures' from 'lindb.storage.replicator.runner' where type='local' group by db,node",
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
                sql: "select 'replica_panics' from 'lindb.storage.replicator.runner' where type='local' group by db,node",
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
            title: "Number of Replica Rows",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'replica_rows' from 'lindb.storage.replica.local' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Replica Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'replica_rows' from 'lindb.storage.replica.local' group by db,node",
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
            title: "Ack Replica",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'ack_sequence' from 'lindb.storage.replica.local' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Invalid Replica Sequence",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'invalid_sequence' from 'lindb.storage.replica.local' group by db,node",
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
