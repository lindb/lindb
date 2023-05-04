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

export const MasterControllerDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Master Failover",
            description: "elected new master node",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'failovers' from 'lindb.master.controller' group by node",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Master Failover Failures",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'failover_failures' from 'lindb.master.controller' group by node",
                watch: ["node"],
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
            title: "Master Reassign",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'reassigns' from 'lindb.master.controller' group by node",
                watch: ["node"],
              },
            ],
            unit: Unit.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Master Reassign Failure",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'reassign_failures' from 'lindb.master.controller' group by node",
                watch: ["node"],
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
