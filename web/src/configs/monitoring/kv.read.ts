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

export const KVStoreReadDashboard: Dashboard = {
  variates: [
    {
      tagKey: "node",
      label: "Node",
      db: MonitoringDB,
      multiple: true,
      sql: "show tag values from 'lindb.kv.table.read' with key=node",
    },
  ],
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Read(QPS)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select rate('get_counts') from 'lindb.kv.table.read' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Read Traffic",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'read_bytes' from 'lindb.kv.table.read' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Read Error",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'get_errors' from 'lindb.kv.table.read' group by node",
                watch: ["node"],
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
            title: "MMap File",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'mmap_counts' from 'lindb.kv.table.read' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "MMap File Error",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'mmap_errors' from 'lindb.kv.table.read' group by node",
                watch: ["node"],
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
            title: "UNMMap File",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'unmmap_counts' from 'lindb.kv.table.read' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "UNMMap File Error",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'unmmap_errors' from 'lindb.kv.table.read' group by node",
                watch: ["node"],
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
            title: "Current Active Reader",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'active_readers' from 'lindb.kv.table.cache' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Evict Reader From Cache",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'evict_counts' from 'lindb.kv.table.cache' group by node",
                watch: ["node"],
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
            title: "Hit Reader Cache",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'cache_hits' from 'lindb.kv.table.cache' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Miss Reader Cache",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'cache_misses' from 'lindb.kv.table.cache' group by node",
                watch: ["node"],
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
            title: "Close Reader Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'close_counts' from 'lindb.kv.table.cache' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Close Reader Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'close_errors' from 'lindb.kv.table.cache' group by node",
                watch: ["node"],
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
