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
import { Node, Unit } from "@src/models";
import { StateKit, FormatKit } from "@src/utils";
import React, { CSSProperties, ReactNode, useContext } from "react";
import { URLStore } from "@src/stores";
import { isEmpty, get, concat } from "lodash-es";
import { UIContext } from "@src/context/UIContextProvider";

const { Text } = Typography;
const { CPU, Memory } = StateMetricName;

interface NodeViewProps {
  title?: string;
  style?: CSSProperties;
  nodes: Node[];
  sql: string;
  showNodeId?: boolean;
  statusTip?: ReactNode;
}
export default function NodeView(props: NodeViewProps) {
  const { showNodeId, title, style, statusTip, nodes, sql } = props;
  const { stateMetric } = useStateMetric(sql);
  const { locale } = useContext(UIContext);
  const { NodeView } = locale;

  const nodeIdCol = {
    title: NodeView.nodeId,
    dataIndex: "id",
    key: "id",
  };

  const columns: any[] = [
    {
      title: NodeView.title,
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
                  key: NodeView.hostIp,
                  value: <Text link>{record.hostIp}</Text>,
                },
                {
                  key: NodeView.hostName,
                  value: <Text link>{record.hostName}</Text>,
                },
                {
                  key: NodeView.httpPort,
                  value: <Text link>{record.httpPort}</Text>,
                },
                {
                  key: NodeView.grpcPort,
                  value: <Text link>{get(record, "grpcPort", "-")}</Text>,
                },
              ]}
            />
          </Space>
        );
      },
    },
    {
      title: NodeView.uptime,
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
            {FormatKit.format(new Date().getTime() - text, Unit.Milliseconds)}
          </Space>
        );
      },
    },
    {
      title: NodeView.version,
      dataIndex: "version",
      key: "version",
    },
    {
      title: NodeView.cpu,
      key: "cpu",
      render: (_text: any, record: any, _index: any) => {
        return (
          <>
            {FormatKit.format(
              100 -
                StateKit.getMetricField(
                  stateMetric,
                  CPU,
                  "idle",
                  `${record.hostIp}:${
                    record.grpcPort ? record.grpcPort : record.httpPort
                  }`
                ) *
                  100,
              Unit.Percent
            )}
          </>
        );
      },
    },
    {
      title: NodeView.memory,
      key: "memory",
      render: (_text: any, record: any, _index: any) => {
        const node = `${record.hostIp}:${
          record.grpcPort ? record.grpcPort : record.httpPort
        }`;
        const total = StateKit.getMetricField(
          stateMetric,
          Memory,
          "total",
          node
        );
        const used = StateKit.getMetricField(stateMetric, Memory, "used", node);
        return (
          <CapacityView
            percent={StateKit.getMetricField(
              stateMetric,
              Memory,
              "usage",
              node
            )}
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

  if (isEmpty(title)) {
    return (
      <Table
        size="small"
        bordered
        columns={showNodeId ? concat([nodeIdCol], columns) : columns}
        dataSource={nodes}
        pagination={false}
        empty={statusTip}
      />
    );
  }
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
        columns={showNodeId ? concat([nodeIdCol], columns) : columns}
        dataSource={nodes}
        pagination={false}
        empty={statusTip}
      />
    </Card>
  );
}
