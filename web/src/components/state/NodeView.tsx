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
import {
  Badge,
  Card,
  Descriptions,
  Space,
  Table,
  Typography,
} from "@douyinfe/semi-ui";
import { IconSettingStroked } from "@douyinfe/semi-icons";
import { CapacityView } from "@src/components";
import { StateMetricName, Route } from "@src/constants";
import { useStateMetric } from "@src/hooks";
import { Node } from "@src/models";
import {
  getMetricField,
  transformMilliseconds,
  transformPercent,
} from "@src/utils";
import React, { CSSProperties } from "react";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";

const { Text } = Typography;
const { CPU, Memory } = StateMetricName;

interface NodeViewProps {
  title: string;
  style?: CSSProperties;
  loading?: boolean;
  nodes: Node[];
  sql: string;
  showNodeId?: boolean;
}
export default function NodeView(props: NodeViewProps) {
  const { showNodeId, title, style, loading, nodes, sql } = props;
  const { stateMetric } = useStateMetric(sql);

  const nodeIdCol = {
    title: "Node Id",
    dataIndex: "id",
    key: "id",
  };

  const columns: any[] = [
    {
      title: "Host Information",
      dataIndex: "hostIp",
      key: "hostIp",
      render: (_text: any, record: Node, _index: any) => {
        return (
          <Space align="center">
            <Descriptions
              className="lin-small-desc"
              row
              size="small"
              data={[
                {
                  key: "Host IP",
                  value: <Text link>{record.hostIp}</Text>,
                },
                {
                  key: "Host Name",
                  value: <Text link>{record.hostName}</Text>,
                },
                {
                  key: "HTTP",
                  value: <Text link>{record.httpPort}</Text>,
                },
                {
                  key: "GRPC",
                  value: <Text link>{record.grpcPort}</Text>,
                },
              ]}
            />
          </Space>
        );
      },
    },
    {
      title: "Uptime",
      dataIndex: "onlineTime",
      key: "onlineTime",
      render: (text: any, _record: any, _index: any) => {
        return (
          <Space align="center">
            <Badge
              dot
              style={{
                width: 12,
                height: 12,
                marginTop: 4,
                marginRight: 4,
                backgroundColor: `var(--semi-color-success)`,
              }}
            />
            {transformMilliseconds(new Date().getTime() - text)}
          </Space>
        );
      },
    },
    {
      title: "Version",
      dataIndex: "version",
      key: "version",
    },
    {
      title: "CPU",
      key: "cpu",
      render: (_text: any, record: any, _index: any) => {
        return (
          <>
            {transformPercent(
              100 -
                getMetricField(
                  stateMetric,
                  CPU,
                  "idle",
                  `${record.hostIp}:${record.grpcPort}`
                ) *
                  100
            )}
          </>
        );
      },
    },
    {
      title: "Memory",
      key: "memory",
      render: (_text: any, record: any, _index: any) => {
        const node = `${record.hostIp}:${record.grpcPort}`;
        const total = getMetricField(stateMetric, Memory, "total", node);
        const used = getMetricField(stateMetric, Memory, "used", node);
        return (
          <CapacityView
            percent={getMetricField(stateMetric, Memory, "usage", node)}
            total={total}
            free={total - used}
            used={used}
          />
        );
      },
    },
    {
      title: "",
      key: "operator",
      render: (_text: any, record: any, _index: any) => (
        <Text link className="lin-link">
          <IconSettingStroked
            onClick={() => {
              URLStore.changeURLParams({
                path: Route.ConfigurationView,
                params: { target: `${record.hostIp}:${record.httpPort}` },
              });
            }}
          />
        </Text>
      ),
    },
  ];

  return (
    <Card
      style={{ ...style }}
      title={title}
      headerStyle={{ padding: 12 }}
      bodyStyle={{ padding: 12 }}
    >
      <Table
        size="small"
        bordered={false}
        columns={showNodeId ? _.concat([nodeIdCol], columns) : columns}
        dataSource={nodes}
        loading={loading}
        pagination={false}
      />
    </Card>
  );
}
