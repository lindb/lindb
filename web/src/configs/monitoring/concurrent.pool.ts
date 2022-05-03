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
import { Dashboard, UnitEnum } from "@src/models";

export const ConcurrentPoolDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Active Wrokers",
            description: "current workers count in use",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select workers_alive from lindb.concurrent.pool group by node,pool_name",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Number Of Workers Created",
            description: "workers created count since start",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select workers_created from lindb.concurrent.pool group by node,pool_name",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Number Of Workers Killed",
            description: "workers killed since start",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select workers_killed from lindb.concurrent.pool group by node,pool_name",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Tasks Consumed",
            description: "task consumed and executed success",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select tasks_consumed from lindb.concurrent.pool group by node,pool_name",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Task Rejected",
            description: "task rejected because pool is busy",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select tasks_rejected from lindb.concurrent.pool group by node,pool_name",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Tasks Execute Panic",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select tasks_panic from lindb.concurrent.pool group by node,pool_name",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Task Waiting Time(P99)",
            config: { type: "area" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from lindb.concurrent.pool.tasks_waiting_duration group by node,pool_name",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 12,
        },
        {
          chart: {
            title: "Task Executing Time(P99)",
            description: "include task waiting time",
            config: { type: "area" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from lindb.concurrent.pool.tasks_executing_duration group by node,pool_name",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 12,
        },
      ],
    },
  ],
};
