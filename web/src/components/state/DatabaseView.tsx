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
import React from "react";
import { Card, Table, Descriptions, Typography } from "@douyinfe/semi-ui";
import * as _ from "lodash-es";
import { Route } from "@src/constants";
import { URLStore } from "@src/stores";

const { Text } = Typography;

interface DatabaseViewProps {
  liveNodes: any;
  storage: any;
  loading: boolean;
}
export default function DatabaseView(props: DatabaseViewProps) {
  const { loading, storage, liveNodes } = props;
  const columns = [
    {
      title: "Name",
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
                  path: Route.DatabaseOverview,
                  params: { db: text, storage: storage.name },
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
      title: "Num. Of Shards",
      dataIndex: "stats.numOfShards",
    },
    {
      title: "Replica Factor",
      dataIndex: "stats.totalReplica",
    },
    {
      title: "Replication Status",
      width: "30%",
      render: (text: any, record: any, index: any) => {
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
  ];
  const getDatabaseList = (): any[] => {
    const rs: any[] = [];
    const databaseMap = _.get(storage, "shardStates", {});
    const databaseNames = _.keys(databaseMap);
    databaseNames.map((name: string) => {
      const db = databaseMap[name];
      const stats = {
        totalReplica: 0,
        availableReplica: 0,
        unavailableReplica: 0,
        numOfShards: 0,
      };
      _.mapValues(db, function (shard: any) {
        const replicas = _.get(shard, "replica.replicas", []);
        stats.numOfShards++;
        stats.totalReplica += replicas.length;
        replicas.map((nodeId: number) => {
          if (_.has(liveNodes, nodeId)) {
            stats.availableReplica++;
          } else {
            stats.unavailableReplica++;
          }
        });
        console.log("shard....", shard, replicas);
      });
      rs.push({ name: name, stats: stats });
    });
    return rs;
  };
  return (
    <>
      <Card
        bordered
        title={"Database List"}
        headerStyle={{ padding: 12 }}
        bodyStyle={{ padding: 12 }}
      >
        <Table
          size="small"
          bordered={false}
          columns={columns}
          dataSource={getDatabaseList()}
          loading={loading}
          pagination={false}
        />
      </Card>
    </>
  );
}
