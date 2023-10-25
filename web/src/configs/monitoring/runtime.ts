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

export const RuntimeDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Sys",
            description: "number of heap bytes obtained from system",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select sys from lindb.runtime.mem group by node",
                watch: ["role", "namespace", "node"],
              },
            ],
            unit: Unit.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Frees",
            description: "number of frees",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select frees from lindb.runtime.mem group by node",
                watch: ["role", "node"],
              },
            ],
            unit: Unit.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Total Alloc",
            description: "bytes allocated even if freed",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select total_alloc from lindb.runtime.mem group by node",
                watch: ["role", "node"],
              },
            ],
            unit: Unit.Bytes,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Heap Alloc",
            description: "bytes allocated and not yet freed",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select heap_alloc from lindb.runtime.mem group by node",
                watch: ["role", "node"],
              },
            ],
            unit: Unit.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Heap Objects",
            description: "total number of allocated objects",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select heap_objects from lindb.runtime.mem group by node",
                watch: ["role", "node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Heap In Used",
            description: "bytes in non-idle span",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select heap_inuse from lindb.runtime.mem group by node",
                watch: ["role", "node"],
              },
            ],
            unit: Unit.Bytes,
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
            description: "Number of goroutines",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select go_goroutines from lindb.runtime group by node",
                watch: ["role", "node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Number of Threads",
            description: "Number of Threads",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select go_threads from lindb.runtime group by node",
                watch: ["role", "node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Next GC Bytes",
            description: "Next GC Bytes",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select next_gc from lindb.runtime.mem group by node",
                watch: ["role", "node"],
              },
            ],
            unit: Unit.Bytes,
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
            description: "GC CPU Fraction",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select gc_cpu_fraction from lindb.runtime group by node",
                watch: ["role", "node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Lookups",
            description: "number of pointer lookups",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select lookups from lindb.runtime.mem group by node",
                watch: ["role", "node"],
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
