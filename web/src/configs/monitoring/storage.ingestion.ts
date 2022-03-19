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

export const StorageIngestionDashboard: Dashboard = {
  variates: [
    {
      tagKey: "db",
      label: "Database",
      db: MonitoringDB,
      sql: "show tag values from 'lindb.tsdb.shard' with key=db",
      watch: { clear: ["node"] },
    },
    {
      tagKey: "node",
      label: "Node",
      watch: { cascade: ["db"] },
      db: MonitoringDB,
      multiple: true,
      sql: "show tag values from 'lindb.tsdb.shard' with key=node",
    },
  ],
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Write Metric Batch",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_batches' from 'lindb.tsdb.shard' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Write Metric",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_metrics' from 'lindb.tsdb.shard' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Write Field Data Points",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_fields' from 'lindb.tsdb.shard' group by db,node",
                watch: ["node", "db"],
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
            title: "Generate Metric ID",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_metric_ids' from 'lindb.tsdb.metadb' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Generate Tag Key ID",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_tag_key_ids' from 'lindb.tsdb.metadb' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Generate Field ID",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_field_ids' from 'lindb.tsdb.metadb' group by db,node",
                watch: ["node", "db"],
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
            title: "Build Inverted Index",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'build_inverted_index_counter' from 'lindb.tsdb.indexdb' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Generate Tag Value Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_tag_value_id_fails' from 'lindb.tsdb.indexdb' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
      ],
    },
  ],
};
