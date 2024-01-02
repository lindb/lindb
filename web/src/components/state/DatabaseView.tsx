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
import React, { useContext } from "react";
import { Card, Table, Descriptions, Typography } from "@douyinfe/semi-ui";
import { get, has, keys, mapValues, orderBy } from "lodash-es";
import { Route } from "@src/constants";
import { URLStore } from "@src/stores";
import { UIContext } from "@src/context/UIContextProvider";

const { Text } = Typography;

interface DatabaseViewProps {
  liveNodes: any;
  storage: any;
  loading?: boolean;
  title?: string;
  databaseName?: string;
}
export default function DatabaseView(props: DatabaseViewProps) {
  const { loading, storage, liveNodes, title, databaseName } = props;
  const { locale } = useContext(UIContext);
  const { DatabaseView, StorageView } = locale;
  const columns = [
    {
      title: DatabaseView.name,
      dataIndex: "name",
      key: "name",
      render: (text: any) => {
        return (
          <Text
            link
            className="lin-link"
            onClick={() => {
              if (!text) {
                // only in storage list can click
                URLStore.changeURLParams({
                  path: Route.MonitoringReplication,
                  params: { db: text },
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
      title: DatabaseView.numOfShards,
      dataIndex: "stats.numOfShards",
    },
    {
      title: DatabaseView.replicaFactor,
      dataIndex: "stats.totalReplica",
    },
    {
      title: StorageView.replicationStatus,
      width: "30%",
      render: (_text: any, record: any, _index: any) => {
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
  ];
  const getDatabaseList = (): any[] => {
    const rs: any[] = [];
    const databaseMap = get(storage, "shardStates", {});
    const databaseNames = keys(databaseMap);
    databaseNames.map((name: string) => {
      if (databaseName && databaseName !== name) {
        // if database selected, only show selected database's state
        return;
      }
      const db = databaseMap[name];
      const stats = {
        totalReplica: 0,
        availableReplica: 0,
        unavailableReplica: 0,
        numOfShards: 0,
      };
      mapValues(db, function (shard: any) {
        const replicas = get(shard, "replica.replicas", []);
        stats.numOfShards++;
        stats.totalReplica += replicas.length;
        replicas.map((nodeId: number) => {
          if (has(liveNodes, nodeId)) {
            stats.availableReplica++;
          } else {
            stats.unavailableReplica++;
          }
        });
      });
      rs.push({ name: name, stats: stats });
    });
    return rs;
  };
  return (
    <>
      <Table
        size="small"
        bordered
        columns={columns}
        dataSource={orderBy(getDatabaseList(), ["name"], ["asc"])}
        loading={loading}
        pagination={false}
      />
    </>
  );
}
