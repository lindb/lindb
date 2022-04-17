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

export const BrokerSystemsDashboard: Dashboard = {
  variates: [
    {
      tagKey: "node",
      label: "Node",
      db: MonitoringDB,
      multiple: true,
      sql: "show tag values from 'lindb.monitor.system.mem_stat' with key=node where role=Broker",
    },
  ],
  rows: [
    {
      panels: [
        {
          chart: {
            title: "CPU Usage",
            config: { type: "area" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 100-idle*100 as used_percent from lindb.monitor.system.cpu_stat where role=Broker group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Percent,
          },
          span: 12,
        },
        {
          chart: {
            title: "Memory Usage",
            config: { type: "area" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select used_percent from lindb.monitor.system.mem_stat where role=Broker group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Percent,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Net In Speed",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select bytes_recv/1024 from lindb.monitor.system.net_stat where role=Broker group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.KBytesPerSecond,
          },
          span: 12,
        },
        {
          chart: {
            title: "Net Out Speed",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select bytes_sent/1024 from lindb.monitor.system.net_stat where role=Broker group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.KBytesPerSecond,
          },
          span: 12,
        },
      ],
    },
  ],
};
