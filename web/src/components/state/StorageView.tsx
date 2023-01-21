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
import React, { useContext } from "react";
import { URLStore } from "@src/stores";
import { UIContext } from "@src/context/UIContextProvider";

const { Text } = Typography;

const StorageView: React.FC<{
  name?: string;
  storages: StorageState[];
  loading?: boolean;
  statusTip?: React.ReactNode;
}> = (props) => {
  const { name, loading, storages, statusTip } = props;
  const { locale } = useContext(UIContext);
  const { StorageView } = locale;

  const columns = [
    {
      title: StorageView.name,
      dataIndex: "name",
      key: "name",
      render: (text: any) => {
        return (
          <Text
            link
            className="lin-link"
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
      title: StorageView.nodeStatus,
      render: (_text: any, record: StorageState, _index: any) => {
        return (
          <Descriptions
            row
            className="lin-small-desc"
            size="small"
            data={[
              {
                key: StorageView.aliveNodes,
                value: (
                  <Text type="success">
                    {_.get(record, "stats.liveNodes", 0)}
                  </Text>
                ),
              },
              {
                key: StorageView.deadNodes,
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
      title: StorageView.numOfDatabase,
      key: "num_db",
      render: (_text: any, record: StorageState, _index: any) => {
        return _.keys(_.get(record, "shardStates", {})).length;
      },
    },
    {
      title: StorageView.replicationStatus,
      width: "30%",
      render: (_text: any, record: StorageState, _index: any) => {
        return (
          <Descriptions
            className="lin-small-desc"
            row
            size="small"
            data={[
              {
                key: StorageView.totalOfReplication,
                value: (
                  <Text link>{_.get(record, "stats.totalReplica", 0)}</Text>
                ),
              },
              {
                key: StorageView.underReplicated,
                value: (
                  <Text type="success">
                    {_.get(record, "stats.availableReplica", 0)}
                  </Text>
                ),
              },
              {
                key: StorageView.unavailableReplica,
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
      title: StorageView.diskCapacityUsage,
      render: (_text: any, record: any) => {
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
      title={name ? "" : StorageView.storageClusterList}
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
        empty={statusTip}
      />
    </Card>
  );
};

export default StorageView;
