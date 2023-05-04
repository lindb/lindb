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

export const TSDBJobDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Number Of Flush Request Inflight",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'flush_inflight' from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
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
            title: "Flush Data Job",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as count from 'lindb.tsdb.shard.memdb_flush_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Flush Data Job Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'memdb_flush_failures' from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Flush Data Duration(P99)",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.tsdb.shard.memdb_flush_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Milliseconds,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Flush Metric Meta Job",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as count from 'lindb.tsdb.database.metadb_flush_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Flush Metric Meta Job Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'metadb_flush_failures' from 'lindb.tsdb.database' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Flush Metric Meta Duration(P99)",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.tsdb.database.metadb_flush_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Milliseconds,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Flush Index Job",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as count from 'lindb.tsdb.shard.indexdb_flush_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Flush Index Job Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'indexdb_flush_failures' from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Flush Index Duration(P99)",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.tsdb.shard.indexdb_flush_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Milliseconds,
          },
          span: 8,
        },
      ],
    },
  ],
};
