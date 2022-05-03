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

export const BrokerQueryDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Current Executing Task(Alive)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select alive_tasks from lindb.broker.query group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Create Query Tasks",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select created_tasks from lindb.broker.query group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Expired Tasks",
            description: "long-term no response",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select expire_tasks from lindb.broker.query group by node",
                watch: ["namespace", "node", "role"],
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
            title: "Emit Response To Parent Node",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select emitted_responses from lindb.broker.query group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Omit Response Because Task Evicted",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select omitted_responses from lindb.broker.query group by node",
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
            title: "Sent Request",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select sent_requests from lindb.broker.query group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Rend Requst Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select sent_requests_failures from lindb.broker.query group by node",
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
            title: "Send Response",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select sent_responses from lindb.broker.query group by node",
                watch: ["namespace", "node", "role"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Send Response Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select sent_responses_failures from lindb.broker.query group by node",
                watch: ["namespace", "node", "role"],
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
