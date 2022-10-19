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

export const KVStoreJobDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Flushing",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'flushing' from 'lindb.kv.flush' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Complete Flush Job",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as flush from 'lindb.kv.flush.duration' group by node",
                watch: ["node", "namespace"],
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
            title: "Flush Job Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'failure' from 'lindb.kv.flush' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Flush Duration(P99)",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.kv.flush.duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Milliseconds,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Compacting",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'compacting' from 'lindb.kv.compaction' group by node,type",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Complete Compact Job",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as compact from 'lindb.kv.compaction.duration' group by node,type",
                watch: ["node", "namespace"],
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
            title: "Compact Job Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'failure' from 'lindb.kv.compaction' group by node,type",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Compact Duration(P99)",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.kv.compaction.duration' group by node,type",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Milliseconds,
          },
          span: 12,
        },
      ],
    },
  ],
};
