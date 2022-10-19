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

export const ConcurrentLimitDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Processed Requests",
            description: "number of processed requests",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select processed from lindb.concurrent.limit group by node,type",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Throttles",
            description: "number of reaches the max-concurrency",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select throttle_requests from lindb.concurrent.limit group by node,type",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Timeouts",
            description: "number pending and then timeout",
            config: { type: "line", options: chartOptions },
            targets: [
              {
                db: MonitoringDB,
                sql: "select timeout_requests from lindb.concurrent.limit group by node,type",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: Unit.Short,
          },
          span: 8,
        },
      ],
    },
  ],
};
