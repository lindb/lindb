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

export const StorageQueryDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Metric Query(Plan)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select metric_queries from lindb.storage.query group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Metric Query(Plan) Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select metric_query_failures from lindb.storage.query group by node",
                watch: ["namespace", "node", "role"],
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
            title: "Meta Query",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select meta_queries from lindb.storage.query group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Meta Query Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select meta_query_failures from lindb.storage.query group by node",
                watch: ["namespace", "node", "role"],
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
            title: "Ommited Rquests",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select omitted_requests from lindb.storage.query group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 24,
        },
      ],
    },
  ],
};
