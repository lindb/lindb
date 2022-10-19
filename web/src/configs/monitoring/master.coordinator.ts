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

export const MasterCoordinatorDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Storage Node Joins",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'handle_events' from 'lindb.master.state_manager' where type='NodeStartup' group by node",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Storage Node Leaves",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'handle_events' from 'lindb.master.state_manager' where type='NodeFailure' group by node",
                watch: ["node"],
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
            title: "Database Config Changed",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'handle_events' from 'lindb.master.state_manager' where type='DatabaseConfigChanged' group by node",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Database Config Deletion",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'handle_events' from 'lindb.master.state_manager' where type='DatabaseConfigDeletion' group by node",
                watch: ["node"],
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
            title: "Storage Config Changed",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'handle_events' from 'lindb.master.state_manager' where type='StorageConfigChanged' group by node",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Storage Config Deletion",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'handle_events' from 'lindb.master.state_manager' where type='StorageConfigDeletion' group by node",
                watch: ["node"],
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
            title: "Shard Leader Election",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'elections' from 'lindb.master.shard.leader' group by node",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Shard Leader Election Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'elect_failures' from 'lindb.master.shard.leader' group by node",
                watch: ["node"],
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
            title: "Shard Assigns",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'handle_events' from 'lindb.master.state_manager' where type='ShardAssignmentChanged' group by node",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Failure(Process Event)",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'handle_event_failures' from 'lindb.master.state_manager' group by node,type",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Panic(Process Event)",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'panics' from 'lindb.master.state_manager' group by node",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
      ],
    },
  ],
};
