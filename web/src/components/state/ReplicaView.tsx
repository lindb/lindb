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
  Button,
  Card,
  List,
  Popover,
  Space,
  Switch,
  Tag,
  Typography,
} from "@douyinfe/semi-ui";
import { ReplicaState } from "@src/models";
import { exec } from "@src/services";
import { getColor } from "@src/utils";
import * as _ from "lodash-es";
import React, { useEffect, useState } from "react";

const { Text } = Typography;

interface ReplicaViewProps {
  db: string;
  storage: string;
}

export default function ReplicaView(props: ReplicaViewProps) {
  const { db, storage } = props;
  const [loading, setLoading] = useState(true);
  const [replicaState, setReplicaState] = useState<ReplicaState>();
  const [showShard, setShowShard] = useState(true);
  const [showLag, setShowLag] = useState(true);
  const buildReplicaState = (): any[] => {
    const nodes = _.keys(replicaState);
    const rs: any[] = [];
    _.forEach(nodes, (node) => {
      const logs = _.get(replicaState, node, []);
      const shards: any[] = [];
      _.forEach(logs, (log) => {
        const shardIdx = _.findIndex(shards, { shardId: log.shardId });
        const families = _.get(log, "replicators", []);
        const totalPending = _.reduce(
          families,
          function (sum, n) {
            return sum + n.pending;
          },
          0
        );
        if (shardIdx < 0) {
          shards.push({
            shardId: log.shardId,
            pending: totalPending,
            families: [log],
          });
        } else {
          const shard = shards[shardIdx];
          shard.pending += totalPending;
          shard.families.push(log);
        }
      });
      rs.push({ node: node, shards: shards });
    });
    return rs;
  };
  useEffect(() => {
    const fetchReplicaState = async (sql: string) => {
      try {
        setLoading(true);
        const state = await exec<ReplicaState>({ sql: sql });
        setReplicaState(state);
      } catch (err) {
        console.log(err);
      } finally {
        setLoading(false);
      }
    };
    fetchReplicaState(
      `show replication where storage='${storage}' and database='${db}'`
    );
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const renderShardDetail = (shard: any) => {
    return (
      <>
        <List
          bordered
          size="small"
          dataSource={_.orderBy(shard.families, ["familyTime"], ["desc"])}
          renderItem={(item) => (
            <List.Item key={item.familyTime}>
              <Space align="center">
                <div style={{ textAlign: "center", rowGap: 4 }}>
                  <div>{item.familyTime}</div>
                  <div style={{ columnGap: 4, display: "flex" }}>
                    <Tag color="blue">Leader:{item.leader}</Tag>
                    <Tag color="blue">
                      Append:{item.append > 0 ? item.append - 1 : 0}
                    </Tag>
                  </div>
                </div>
                <div>
                  {_.get(item, "replicators", []).map((r: any) => (
                    <>
                      <div
                        style={{
                          border: "1px solid var(--semi-color-border)",
                          marginBottom: 4,
                          padding: 4,
                          display: "flex",
                          columnGap: 4,
                        }}
                      >
                        <Tag color="blue">
                          Peer:{item.leader}=&gt;{r.replicator}
                        </Tag>
                        <Tag color="blue">
                          Consume:{r.consume > 0 ? r.consume - 1 : 0}
                        </Tag>
                        <Tag color="blue">Ack:{r.ack}</Tag>
                        <Tag color="blue">Lag:{r.pending}</Tag>
                      </div>
                    </>
                  ))}
                </div>
              </Space>
            </List.Item>
          )}
        />
      </>
    );
  };
  return (
    <Card
      bordered
      title="Replication Status"
      headerExtraContent={
        <>
          <div style={{ display: "flex", alignItems: "center" }}>
            <Text style={{ margin: 4 }}>Shard</Text>
            <Switch checked={showShard} onChange={setShowShard} size="small" />
            <Text style={{ margin: 4 }}>Lag</Text>
            <Switch checked={showLag} onChange={setShowLag} size="small" />
          </div>
        </>
      }
      headerStyle={{ padding: 12 }}
      bodyStyle={{ padding: "12px 12px 0px 12px", display: "flex" }}
      loading={loading}
    >
      <Space wrap>
        {buildReplicaState().map((item) => {
          return (
            <Card
              bordered
              style={{ margin: "0px 10px 10px 0px", width: 300 }}
              bodyStyle={{ padding: 12 }}
              key={item.node}
            >
              <div style={{ marginBottom: 2, textAlign: "center" }}>
                {item.node}
              </div>
              {_.get(item, "shards", []).map((shard: any) => {
                return (
                  <Popover
                    showArrow
                    key={shard.shardId}
                    content={renderShardDetail(shard)}
                  >
                    <Button
                      style={{
                        minWidth: 50,
                        margin: "0px 4px 4px 0px",
                        color: "var(--semi-color-text-0)",
                        backgroundColor: getColor(shard.shardId),
                      }}
                      size="small"
                    >
                      {showShard ? <span>S:{shard.shardId}</span> : ""}
                      {showLag ? <span>L:{shard.pending}</span> : ""}
                    </Button>
                  </Popover>
                );
              })}
            </Card>
          );
        })}
      </Space>
    </Card>
  );
}
