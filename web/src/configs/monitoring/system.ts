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
import { Dashboard, LegendAggregateType, Unit } from "@src/models";

export const chartOptions = {
  legend: {
    values: [
      LegendAggregateType.MAX,
      LegendAggregateType.MIN,
      LegendAggregateType.CURRENT,
    ],
  },
};

export const SystemDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "CPU Usage",
            config: {
              type: "area",
              options: chartOptions,
            },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 100-idle*100 as usage from lindb.monitor.system.cpu_stat group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Percent,
          },
          span: 12,
        },
        {
          chart: {
            title: "CPU",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select idle,system,iowait,user from lindb.monitor.system.cpu_stat group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Percent2,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Memory Usage",
            config: { type: "area", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select usage from lindb.monitor.system.mem_stat group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Percent,
          },
          span: 12,
        },
        {
          chart: {
            title: "Memory Size",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select total,free,used from lindb.monitor.system.mem_stat group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Bytes,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Disk Usage",
            config: { type: "area", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select usage from lindb.monitor.system.disk_usage_stats group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Percent,
          },
          span: 8,
        },
        {
          chart: {
            title: "INodes Usage",
            config: { type: "area", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select inodes_usage from lindb.monitor.system.disk_inodes_stats group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Percent,
          },
          span: 8,
        },
        {
          chart: {
            title: "Disk Size",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select total,used,free from lindb.monitor.system.disk_usage_stats group by node",
                watch: ["namespace", "node", "role"],
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
            title: "Number Of Bytes Sent",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select bytes_sent/1024 as sent from lindb.monitor.system.net_stat group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.KBytesPerSecond,
          },
          span: 12,
        },
        {
          chart: {
            title: "Number Of Bytes Receive",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select bytes_recv/1024 as receive from lindb.monitor.system.net_stat group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.KBytesPerSecond,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Number Of Packages Sent",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select packets_sent as sent from lindb.monitor.system.net_stat group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Number Of Packages Receive",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select packets_recv as receive from lindb.monitor.system.net_stat group by node",
                watch: ["namespace", "node", "role"],
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
            title: "Number Of Errors While Sending",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select errin from lindb.monitor.system.net_stat group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Number Of Errors While Receiving",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select errin from lindb.monitor.system.net_stat group by node",
                watch: ["namespace", "node", "role"],
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
