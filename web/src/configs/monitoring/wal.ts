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

export const WALDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Write Traffic",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'receive_write_bytes' from 'lindb.storage.wal' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Write",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_wal' from 'lindb.storage.wal' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Write Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'write_wal_failures' from 'lindb.storage.wal' group by node",
                watch: ["node", "namespace"],
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
            title: "Replica Traffic",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'receive_replica_bytes' from 'lindb.storage.wal' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Bytes,
          },
          span: 8,
        },
        {
          chart: {
            title: "Replica",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'replica_wal' from 'lindb.storage.wal' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Replica Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'replica_wal_failures' from 'lindb.storage.wal' group by node",
                watch: ["node", "namespace"],
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
