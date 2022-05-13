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
import { Card, Descriptions, Table, Typography } from "@douyinfe/semi-ui";
import { DiskUsageView } from "@src/components";
import { Route, StateMetricName } from "@src/constants";
import { StorageState } from "@src/models";
import * as _ from "lodash-es";
import React from "react";
import { URLStore } from "@src/stores";

const { Text } = Typography;

interface StorageViewProps {
  name?: string;
  storages: StorageState[];
  loading?: boolean;
}

export default function StorageView(props: StorageViewProps) {
  const { name, loading, storages } = props;

  const columns = [
    {
      title: "Name(Namespace)",
      dataIndex: "name",
      key: "name",
      render: (text: any) => {
        return (
          <Text
            link
            onClick={() => {
              if (!name) {
                // only in storage list can click
                URLStore.changeURLParams({
                  path: Route.StorageOverview,
                  params: { name: text },
                });
              }
            }}
          >
            {text}
          </Text>
        );
      },
    },
    {
      title: "Node Status",
      render: (text: any, record: StorageState, index: any) => {
        return (
          <Descriptions
            row
            className="lin-small-desc"
            size="small"
            data={[
              {
                key: "Alive Nodes",
                value: (
                  <Text type="success">
                    {_.get(record, "stats.liveNodes", 0)}
                  </Text>
                ),
              },
              {
                key: "Dead Nodes",
                value: (
                  <Text type="danger">
                    {_.get(record, "stats.deadNodes.length", 0)}
                  </Text>
                ),
              },
            ]}
          />
        );
      },
    },
    {
      title: "Num. Of Database",
      key: "num_db",
      render: (text: any, record: StorageState, index: any) => {
        return _.keys(_.get(record, "shardStates", {})).length;
      },
    },
    {
      title: "Replication Status",
      width: "30%",
      render: (text: any, record: StorageState, index: any) => {
        return (
          <Descriptions
            className="lin-small-desc"
            row
            size="small"
            data={[
              {
                key: "Total",
                value: (
                  <Text link>{_.get(record, "stats.totalReplica", 0)}</Text>
                ),
              },
              {
                key: "Under-replicated",
                value: (
                  <Text type="success">
                    {_.get(record, "stats.availableReplica", 0)}
                  </Text>
                ),
              },
              {
                key: "Unavailable",
                value: (
                  <Text type="danger">
                    {_.get(record, "stats.unavailableReplica", 0)}
                  </Text>
                ),
              },
            ]}
          />
        );
      },
    },
    {
      title: "Disk Capacity Usage",
      render: (text: any, record: any, index: any) => {
        return (
          <DiskUsageView
            sql={`show storage metric where storage='${record.name}' and metric in ('${StateMetricName.Disk}')`}
          />
        );
      },
    },
  ];

  return (
    <Card
      title={name ? "" : "Storage Cluster List"}
      headerStyle={{ padding: 12 }}
      bodyStyle={{ padding: 12 }}
    >
      <Table
        size="small"
        bordered={false}
        columns={columns}
        dataSource={storages}
        loading={loading}
        pagination={false}
      />
    </Card>
  );
}
