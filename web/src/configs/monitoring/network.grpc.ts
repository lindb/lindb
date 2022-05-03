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

export const NetworkGRPCDasbhoard: Dashboard = {
  rows: [
    {
      panels: [
        {
          chart: {
            title: "Number Of Send(Client Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as sent from 'lindb.traffic.grpc_client.stream.sent_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Send Failure(Client Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'msg_sent_failures' as failure from 'lindb.traffic.grpc_client.stream' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Send Duration(P99)(Client Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.traffic.grpc_client.stream.sent_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Number Of Receive(Client Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as received from 'lindb.traffic.grpc_client.stream.received_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Receive Failure(Client Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'msg_received_failures' as failure from 'lindb.traffic.grpc_client.stream' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Receive Duration(P99)(Client Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.traffic.grpc_client.stream.received_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Number Of Send(Server Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as sent from 'lindb.traffic.grpc_server.stream.sent_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Send Failure(Server Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'msg_sent_failures' as failure from 'lindb.traffic.grpc_server.stream' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Send Duration(P99)(Server Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.traffic.grpc_server.stream.sent_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Number Of Receive(Server Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as received from 'lindb.traffic.grpc_server.stream.received_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Receive Failure(Server Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'msg_received_failures' as failure from 'lindb.traffic.grpc_server.stream' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Receive Duration(P99)(Server Stream)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.traffic.grpc_server.stream.received_duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Handle Request(Client Unary)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as hander from 'lindb.traffic.grpc_client.unary.duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Handle Request Failure(Client Unary)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'failures' from 'lindb.traffic.grpc_client.unary' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Handle Request Duration(P99)(Client Unary)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.traffic.grpc_client.unary.duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Handle Request(Server Unary)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'HistogramCount' as hander from 'lindb.traffic.grpc_server.unary.duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Handle Request Failure(Server Unary)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'failures' from 'lindb.traffic.grpc_server.unary' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Short,
          },
          span: 8,
        },
        {
          chart: {
            title: "Handle Request Duration(P99)(Server Unary)",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select quantile(0.99) as p99 from 'lindb.traffic.grpc_server.unary.duration' group by node",
                watch: ["node", "namespace"],
              },
            ],
            unit: UnitEnum.Milliseconds,
          },
          span: 8,
        },
      ],
    },
    {
      panels: [
        {
          chart: {
            title: "Server Handle Panic",
            config: { type: "line" },
            targets: [
              {
                db: MonitoringDB,
                sql: "select 'panics' from 'lindb.traffic.grpc_server' group by node",
                watch: ["node", "namespace"],
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
