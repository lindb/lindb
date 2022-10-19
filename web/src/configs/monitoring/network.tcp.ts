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

export const NetworkTCPDasbhoard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Current Active Connections",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'active_conns' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Accpet",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'accept_conns' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Accpet Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'accept_failures' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Number Of Read",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'reads' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Read Traffic",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'read_bytes' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Read Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'read_failures' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Number Of Write",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'writes' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "WRite Traffic",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_bytes' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Write Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_failures' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Closed Connections",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'close_conns' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Close Connection Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'close_failures' from 'lindb.traffic.tcp' group by node",
                watch: ["node", "namespace"],
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
