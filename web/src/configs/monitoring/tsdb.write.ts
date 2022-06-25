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

export const TSDBWriteDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Current Active Families",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select sum('active_families') as families from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Current Active Memory Databases",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select sum('active_memdbs') as memdb from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Memory Databases Total Size",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select sum('memdb_total_size') as memdb_size from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Bytes,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Batch",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_batches' from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Write Failures",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_metrics_failures' from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Write Metrics",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_metrics' from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Write Fields",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_fields' from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
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
                sql: "select 'build_inverted_index' from 'lindb.tsdb.indexdb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Lookup Metric Meta Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'lookup_metric_meta_failures' from 'lindb.tsdb.shard' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Generate Metric Id",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_metric_ids' from 'lindb.tsdb.metadb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Generate Metric Id Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_metric_id_failures' from 'lindb.tsdb.metadb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Generate Field Id",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_field_ids' from 'lindb.tsdb.metadb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Generate Field Id Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_field_id_failures' from 'lindb.tsdb.metadb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Generate Tag Key Id",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_tag_key_ids' from 'lindb.tsdb.metadb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Generate Tag Key Id Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_tag_key_id_failures' from 'lindb.tsdb.metadb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Generate Tag Value Id",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_tag_value_ids' from 'lindb.tsdb.metadb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Generate Tag Value Id Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'gen_tag_value_id_failures' from 'lindb.tsdb.metadb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Allocated Pages",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'allocated_pages' from 'lindb.tsdb.memdb' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Allocate Page Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'allocated_page_failures' from 'lindb.tsdb.memdb' group by node",
                watch: ["node", "namespace"],
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
