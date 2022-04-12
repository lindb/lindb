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

export const BrokerIngestionDashboard: Dashboard = {
  variates: [
    {
      tagKey: "db",
      label: "Database",
      db: MonitoringDB,
      sql: "show tag values from 'lindb.broker.replica' with key=db",
      watch: { clear: ["node"] },
    },
    {
      tagKey: "node",
      label: "Node",
      watch: { cascade: ["db"] },
      db: MonitoringDB,
      multiple: true,
      sql: "show tag values from 'lindb.broker.replica' with key=node",
    },
  ],
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Current Active Family Channels",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'active_families' from 'lindb.broker.replica' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Batch Metric",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select rate('batch_metrics') from 'lindb.broker.replica' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Batch Metric Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select rate('batch_metrics_failures') from 'lindb.broker.replica' group by db,node",
                watch: ["node", "db"],
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
            title: "Sent Successfully",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select rate('send_success') from 'lindb.broker.replica' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Sent Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select rate('send_failure') from 'lindb.broker.replica' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Sent Size",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'send_size' from 'lindb.broker.replica' group by db,node",
                watch: ["node", "db"],
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
            title: "Pending For Sending",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'pending_send' from 'lindb.broker.replica' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Retry Send",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'retry' from 'lindb.broker.replica' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Drop After Retry Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'retry_drop' from 'lindb.broker.replica' group by db,node",
                watch: ["node", "db"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
      ],
    },
  ],
};
