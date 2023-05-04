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
import React, { useContext, useEffect, useState } from "react";
import {
  Select,
  Avatar,
  Card,
  List,
  Table,
  Typography,
  Button,
} from "@douyinfe/semi-ui";
import { IconDoubleChevronRight } from "@douyinfe/semi-icons";
import { MemoryDatabaseState, Unit } from "@src/models";
import { Icon } from "@src/components";
import * as _ from "lodash-es";
import { useParams } from "@src/hooks";
import { ColorKit, FormatKit } from "@src/utils";
import { URLStore } from "@src/stores";
import { UIContext } from "@src/context/UIContextProvider";

const MemoryFilterForm: React.FC<{ state: any }> = (props) => {
  const { state } = props;
  const { shard, family, show } = useParams(["shard", "family", "show"]);
  const { locale } = useContext(UIContext);
  const { ReplicationView } = locale;

  return (
    <>
      <>
        <Select
          insetLabel={ReplicationView.shard}
          style={{ marginRight: 8 }}
          showClear
          value={shard}
          onClear={() => {
            URLStore.changeURLParams({ needDelete: ["shard"] });
          }}
          optionList={_.orderBy(
            _.map(state.shards, (shard: any) => {
              return { label: shard.shardId, value: shard.shardId };
            }),
            ["label"]
          )}
          onChange={(v) => {
            URLStore.changeURLParams({ params: { shard: v } });
          }}
        />
        <Select
          insetLabel="Family"
          showClear
          value={family}
          onClear={() => {
            URLStore.changeURLParams({ needDelete: ["family"] });
          }}
          optionList={_.orderBy(
            _.map(state.families, (f: any) => {
              return { label: f, value: f };
            }),
            ["label"],
            ["desc"]
          )}
          onChange={(v) => {
            URLStore.changeURLParams({ params: { family: v } });
          }}
        />
        <Button
          style={{ marginLeft: 8 }}
          icon={
            show === "replica" ? (
              <Icon icon="iconbx-git-repo-forked" style={{ fontSize: 14 }} />
            ) : (
              <Icon icon="icondatabase" style={{ fontSize: 14 }} />
            )
          }
          onClick={() => {
            URLStore.changeURLParams({
              params: { show: show === "replica" ? "memory" : "replica" },
            });
          }}
        />
      </>
    </>
  );
};

const MemoryDatabaseView: React.FC<{
  liveNodes: Node[];
  state: MemoryDatabaseState;
}> = (props) => {
  const { liveNodes, state } = props;
  const { shard, family, show } = useParams(["shard", "family", "show"]);
  const [memoryState, setMemoryState] = useState<any>({
    shards: [],
    families: [],
  });
  const { locale } = useContext(UIContext);
  const { ReplicationView } = locale;

  useEffect(() => {
    const rs: any[] = [];
    const nodes = _.keys(state);
    const shardMap = new Map();
    const familyMap = new Map();
    _.forEach(nodes, (node) => {
      const familyList = _.get(state, node, []);
      _.forEach(familyList, (family) => {
        const shardId = family.shardId;
        const familyTime = family.familyTime;
        familyMap.set(familyTime, familyTime);

        const databases = family.memoryDatabases || [];
        const replica: any[] = [];
        _.forIn(family.replicaSequences, (v, k) => {
          replica.push({ nodeId: k, replicaSeq: v });
        });
        const ack = family.ackSequences;
        if (shardMap.has(shardId)) {
          const shard = shardMap.get(shardId);
          const channels = _.find(shard.channels, { familyTime: familyTime });
          if (channels) {
            // add family into exist family channel
            channels.databases.push(...databases);
          } else {
            // add family into exist shard as new channel
            shard.channels.push({
              node: node,
              familyTime: familyTime,
              databases: databases,
              replica: replica,
              ack: ack,
            });
          }
        } else {
          const shard = {
            shardId: shardId,
            channels: [
              {
                node: node,
                familyTime: familyTime,
                databases: databases,
                replica: replica,
                ack: ack,
              },
            ],
          };
          shardMap.set(shardId, shard);
          rs.push(shard);
        }
      });
    });
    const shards = _.orderBy(rs, ["shardId"]);
    _.forEach(shards, (s) => {
      s.channels = _.orderBy(s.channels, ["familyTime"], ["desc"]);
    });
    setMemoryState({
      shards: shards,
      families: Array.from(familyMap.keys()),
    });
  }, [state]);

  const getNode = (id: string) => {
    const follower = _.find(liveNodes, {
      id: parseInt(id),
    });
    return `${_.get(follower, "hostIp", "unkonw")}:${_.get(
      follower,
      "grpcPort",
      "unkonw"
    )}`;
  };

  const renderSequences = (
    idx: number,
    shardIdx: number,
    node: string,
    replica: any[],
    ack: object
  ) => {
    return (
      <Table
        size="small"
        pagination={false}
        dataSource={replica || []}
        showHeader={idx === 0 && shardIdx === 0}
        columns={[
          {
            title: ReplicationView.node,
            dataIndex: "nodeId",
            render: (text, _record, _index) => {
              return (
                <>
                  <div
                    style={{
                      display: "flex",
                      alignItems: "center",
                    }}
                  >
                    <span>{getNode(text)}</span>
                    <IconDoubleChevronRight
                      style={{
                        marginLeft: 4,
                        color: "var(--semi-color-success)",
                      }}
                    />
                    <span>{node}</span>
                  </div>
                </>
              );
            },
          },
          {
            title: ReplicationView.replica,
            dataIndex: "replicaSeq",
            width: 180,
            render: (text, _record, _index) => {
              return FormatKit.format(text, Unit.Short);
            },
          },
          {
            title: ReplicationView.ack,
            dataIndex: "ack",
            width: 180,
            render: (_text, record, _index) => {
              return FormatKit.format(
                _.get(ack, `${record.nodeId}`) || 0,
                Unit.Short
              );
            },
          },
        ]}
      />
    );
  };

  const renderMemoryDatabase = (
    idx: number,
    shardIdx: number,
    node: string,
    databases: any
  ) => {
    return (
      <Table
        size="small"
        pagination={false}
        dataSource={databases || []}
        showHeader={idx === 0 && shardIdx === 0}
        empty={ReplicationView.noMemoryDatabase}
        columns={[
          {
            title: ReplicationView.node,
            dataIndex: "node",
            width: 160,
            render: (_text, _record, _index) => {
              return node;
            },
          },
          {
            title: ReplicationView.state,
            dataIndex: "state",
            render: (text, _record, _index) => {
              return _.upperFirst(text);
            },
          },
          {
            title: ReplicationView.uptime,
            dataIndex: "uptime",
            width: 100,
            render: (text, _record, _index) => {
              return FormatKit.format(text, Unit.Nanoseconds);
            },
          },
          {
            title: ReplicationView.memSize,
            dataIndex: "memSize",
            width: 130,
            render: (text, _record, _index) => {
              return FormatKit.format(text, Unit.Bytes);
            },
          },
          {
            title: ReplicationView.numOfMetrics,
            dataIndex: "numOfMetrics",
            width: 150,
            render: (text, _record, _index) => {
              return FormatKit.format(text, Unit.Short);
            },
          },
          {
            title: ReplicationView.numOfSeries,
            dataIndex: "numOfSeries",
            width: 150,
            render: (text, _record, _index) => {
              return FormatKit.format(text, Unit.Short);
            },
          },
        ]}
      />
    );
  };

  const renderChannel = (shard: any, shardIdx: any) => {
    return (shard.channelList || []).map((channel: any, idx: any) => {
      return (
        <div
          key={idx}
          style={{
            display: "flex",
          }}
        >
          <Typography.Text
            style={{
              borderLeft: "1px solid var(--semi-color-border)",
              borderRight: "1px solid var(--semi-color-border)",
              borderBottom: "1px solid var(--semi-color-border)",
              minWidth: 180,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            {channel.familyTime}
          </Typography.Text>
          {show === "replica"
            ? renderSequences(
                idx,
                shardIdx,
                channel.node,
                channel.replica,
                channel.ack
              )
            : renderMemoryDatabase(
                idx,
                shardIdx,
                channel.node,
                channel.databases
              )}
        </div>
      );
    });
  };

  const renderShardsAsTable = () => {
    let shards = memoryState.shards;
    if (!_.isEmpty(shard)) {
      shards = _.filter(shards, function (o) {
        return `${o.shardId}` == shard;
      });
    }
    if (!_.isEmpty(family)) {
      shards = _.filter(shards, function (o) {
        o.channelList = _.filter(o.channels, function (f) {
          return f.familyTime == family;
        });
        return !_.isEmpty(o.channelList);
      });
    } else {
      _.forEach(shards, (s) => {
        s.channelList = s.channels;
      });
    }

    if (show !== "replica") {
      // filter empty memory database when show memory database write state
      shards = _.filter(shards, function (o) {
        o.channelList = _.filter(o.channels, function (f) {
          return !_.isEmpty(f.databases);
        });
        return !_.isEmpty(o.channelList);
      });
    }

    return (
      <List
        dataSource={shards}
        renderItem={(item: any, shardIdx: any) => (
          <List.Item
            style={{ display: "block", padding: "0", borderBottom: 0 }}
          >
            <div style={{ display: "flex" }}>
              <div
                style={{
                  borderLeft: "1px solid var(--semi-color-border)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  padding: "0 8px 0 8px",
                  borderBottom: "1px solid var(--semi-color-border)",
                }}
              >
                <Avatar
                  size="small"
                  style={{
                    backgroundColor: ColorKit.getShardColor(item.shardId),
                  }}
                >
                  {item.shardId}
                </Avatar>
              </div>
              <div
                style={{
                  width: "100%",
                  borderRight: "1px solid var(--semi-color-border)",
                }}
              >
                {renderChannel(item, shardIdx)}
              </div>
            </div>
          </List.Item>
        )}
      />
    );
  };

  return (
    <Card
      title={ReplicationView.memoryDatabaseStatus}
      headerStyle={{ padding: 12 }}
      bodyStyle={{ padding: 0 }}
      headerExtraContent={
        <>
          <div style={{ display: "flex", alignItems: "center" }}>
            <MemoryFilterForm state={memoryState} />
          </div>
        </>
      }
    >
      {renderShardsAsTable()}
    </Card>
  );
};

export default MemoryDatabaseView;
