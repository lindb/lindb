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
import { DatabaseView, DiskUsageView, NodeView } from "@src/components";
import { StateMetricName } from "@src/constants";
import { StorageState } from "@src/models";
import { get, keys, orderBy, values } from "lodash-es";
import React, { useContext } from "react";
import { UIContext } from "@src/context/UIContextProvider";

const { Text } = Typography;

const StorageView: React.FC<{
  storage: StorageState;
  loading?: boolean;
  statusTip?: React.ReactNode;
}> = (props) => {
  const { loading, storage, statusTip } = props;
  const { locale } = useContext(UIContext);
  const { StorageView } = locale;

  const columns = [
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
                    {get(record, "stats.liveNodes", 0)}
                  </Text>
                ),
              },
              {
                key: StorageView.deadNodes,
                value: (
                  <Text type="danger">
                    {get(record, "stats.deadNodes.length", 0)}
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
        return keys(get(record, "shardStates", {})).length;
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
                value: <Text link>{get(record, "stats.totalReplica", 0)}</Text>,
              },
              {
                key: StorageView.underReplicated,
                value: (
                  <Text type="success">
                    {get(record, "stats.availableReplica", 0)}
                  </Text>
                ),
              },
              {
                key: StorageView.unavailableReplica,
                value: (
                  <Text type="danger">
                    {get(record, "stats.unavailableReplica", 0)}
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
      render: (_text: any, _: any) => {
        return (
          <DiskUsageView
            sql={`show storage metric where metric in ('${StateMetricName.Disk}')`}
          />
        );
      },
    },
  ];

  return (
    <Card
      title={StorageView.storageCluster}
      headerStyle={{ padding: 12 }}
      bodyStyle={{
        padding: 12,
        display: "flex",
        flexDirection: "column",
        gap: 12,
      }}
    >
      <Table
        size="small"
        bordered
        columns={columns}
        dataSource={[storage]}
        loading={loading}
        pagination={false}
        empty={statusTip}
      />
      <NodeView
        showNodeId
        nodes={orderBy(values(get(storage, "liveNodes", {})), ["id"], ["asc"])}
        sql={`show storage metric where metric in ('${StateMetricName.CPU}','${StateMetricName.Memory}')`}
        style={{ marginTop: 12, marginBottom: 12 }}
      />
      <DatabaseView
        title={StorageView.databaseList}
        liveNodes={get(storage, "liveNodes", {})}
        storage={storage}
      />
    </Card>
  );
};

export default StorageView;
