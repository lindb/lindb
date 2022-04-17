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

export const StorageRuntimeDashboard: Dashboard = {
  variates: [
    {
      tagKey: "db",
      label: "Database",
      db: MonitoringDB,
      sql: "show tag values from 'lindb.runtime.mem' with key=namespace where role=Storage",
      watch: { clear: ["node"] },
    },
    {
      tagKey: "node",
      label: "Node",
      watch: { cascade: ["namespace"] },
      db: MonitoringDB,
      multiple: true,
      sql: "show tag values from 'lindb.runtime.mem' with key=node where role=Storage",
    },
  ],
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Sys (number of heap bytes obtained from system)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select heap_sys_bytes from lindb.runtime.mem where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Frees (number of frees)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select frees_total from lindb.runtime.mem where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Total Alloc (bytes allocated even if freed)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select alloc_bytes_total from lindb.runtime.mem where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
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
            title: "HeapAlloc (bytes allocated and not yet freed)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select heap_alloc_bytes from lindb.runtime.mem where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Heap Objects (total number of allocated objects)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select heap_objects from lindb.runtime.mem where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "HeapInUsed (bytes in non-idle span)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select heap_inuse_bytes from lindb.runtime.mem where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
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
            title: "Number of goroutines",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select go_goroutines from lindb.runtime where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Number of Threads",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select go_threads from lindb.runtime where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Next GC Bytes",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select next_gc_bytes from lindb.runtime.mem where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
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
            title: "GC CPU Fraction",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select gc_cpu_fraction from lindb.runtime where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Lookups(number of pointer lookups)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select lookups_total from lindb.runtime.mem where role=Storage group by namespace,node",
                watch: ["namespace", "node"],
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
