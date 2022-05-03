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

export const IngestionDashboard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Number Of Ingest",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as count from 'lindb.http.ingest_duration' group by node,path",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Ingest Duration(P99)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.http.ingest_duration' group by node,path",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Ingest Metric(Flat)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'ingested_metrics' from 'lindb.ingestion.flat' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Ingest Flat Traffic",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select read_bytes from 'lindb.ingestion.flat' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Bytes,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Data Corrupted(Flat)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'data_corrupted' from 'lindb.ingestion.flat' group by node",
                watch: ["node"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Drop Metrics(Flat)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'dropped_metrics' from 'lindb.ingestion.flat' group by node",
                watch: ["node", "namespace"],
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
            title: "Ingest Metric(Proto)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'ingested_metrics' from 'lindb.ingestion.proto' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Ingest Proto Traffic",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'read_bytes' from 'lindb.ingestion.proto' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Bytes,
          },
          span: 12,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Data Corrupted(Proto)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'data_corrupted' from 'lindb.ingestion.proto' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 12,
        },
        {
          chart: {
            title: "Drop Metric(Proto)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'dropped_metrics' from 'lindb.ingestion.proto' group by node",
                watch: ["node", "namespace"],
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
            title: "Ingest Metric(Influx)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'ingested_metrics' from 'lindb.ingestion.influx' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Ingest Field(Field)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'ingested_fields' from 'lindb.ingestion.influx' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Ingest Traffic(Proto)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'read_bytes' from 'lindb.ingestion.influx' group by node",
                watch: ["node", "namespace"],
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
            title: "Data Corrupted(Influx)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'data_corrupted' from 'lindb.ingestion.influx' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Drop Metric(Influx)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'dropped_metrics' from 'lindb.ingestion.influx' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Drop Field(Influx)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'dropped_fields' from 'lindb.ingestion.influx' group by node",
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
