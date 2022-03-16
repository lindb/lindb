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

export const SystemDashboard: Dashboard = {
  variates: [
    {
      tagKey: "role",
      label: "Role",
      db: MonitoringDB,
      sql: "show tag values from 'lindb.monitor.system.cpu_stat' with key=role",
      watch: { clear: ["node"] },
    },
    {
      tagKey: "node",
      label: "Node",
      watch: { cascade: ["role"] },
      db: MonitoringDB,
      multiple: true,
      sql: "show tag values from 'lindb.monitor.system.cpu_stat' with key=node",
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
                sql: "select 100-idle*100 as used_percent from lindb.monitor.system.cpu_stat group by node",
                watch: ["node", "role"],
              },
            ],
            unit: UnitEnum.Percent,
          },
          span: 12,
        },
        {
          chart: {
            title: "Memory Usage",
            targets: [
              {
                db: MonitoringDB,
                sql: "select used_percent from lindb.monitor.system.mem_stat group by node",
                watch: ["node", "role"],
              },
            ],
            unit: UnitEnum.Percent,
          },
          span: 12,
        },
      ],
    },
  ],
};
