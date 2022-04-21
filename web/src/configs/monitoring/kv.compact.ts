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

export const KVStoreCompactDashboard: Dashboard = {
  variates: [
    {
      tagKey: "node",
      label: "Node",
      db: MonitoringDB,
      multiple: true,
      sql: "show tag values from 'lindb.kv.compaction' with key=node",
    },
  ],
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Compacting",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'compacting' from 'lindb.kv.compaction' group by node,type",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Complete Compact",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as compact from 'lindb.kv.compaction.duration' group by node,type",
                watch: ["node"],
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
            title: "Compact Failure",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'failure' from 'lindb.kv.compaction' group by node,type",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Compact Duration",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) from 'lindb.kv.compaction.duration' group by node,type",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 12,
        },
      ],
    },
  ],
};
