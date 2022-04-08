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
import { Tree, Typography } from "@douyinfe/semi-ui";
import {
  IconFlowChartStroked,
  IconServerStroked,
  IconFile,
} from "@douyinfe/semi-icons";
import { UnitEnum, ChartStatus, ExplainResult } from "@src/models";
import { formatter } from "@src/utils";
import React, { useEffect, useState } from "react";
import { ChartStore } from "@src/stores";
import { reaction } from "mobx";
import { CanvasChart } from "@src/components";
const Text = Typography.Text;

interface ExplainStatsViewProps {
  chartId: string;
}

const ExplainStatsView: React.FC<ExplainStatsViewProps> = (
  props: ExplainStatsViewProps
) => {
  const { chartId } = props;
  const [state, setState] = useState<ExplainResult>();

  const renderCost = (cost: any, total: any) => {
    const percent = (cost * 100.0) / total;
    let type: any = "success";
    if (percent > 50) {
      type = "danger";
    } else if (percent > 30) {
      type = "warning";
    }
    return <Text type={type}>{formatter(cost, UnitEnum.Nanoseconds)}</Text>;
  };

  const buildCollectValueStats = (total: any, key: string, stats: any) => {
    let children: any = [];
    const loadNodes = {
      label: <Text strong>Collect Tag Values(Async)</Text>,
      key: `${key}collect tag values`,
      children: children,
    };
    for (let groupgTagKey of Object.keys(stats)) {
      const cost = stats[groupgTagKey];
      children.push({
        label: (
          <span>
            <Text link strong>
              {groupgTagKey}
            </Text>
            : {renderCost(cost, total)}
          </span>
        ),
        key: `${key + groupgTagKey}cost`,
        // icon: <HddOutlined />,
      });
    }
    return loadNodes;
  };

  const buildLoadStats = (total: any, key: string, loadStats: any) => {
    let children: any = [];
    const loadNodes = {
      label: <Text strong>Load Data(Async)</Text>,
      key: `${key}load data`,
      children: children,
    };
    for (let id of Object.keys(loadStats)) {
      const stats = loadStats[id];
      children.push({
        label: (
          <span>
            <Text strong link>
              {id}
            </Text>
            : [ Count: <Text link>{formatter(stats.count, UnitEnum.None)}</Text>
            , Total Cost: {renderCost(stats.totalCost, total)}, Min:{" "}
            {renderCost(stats.min, total)}, Max: {renderCost(stats.min, total)},
            Num. Of Series:{" "}
            <Text link>{formatter(stats.series, UnitEnum.None)}</Text> ]
          </span>
        ),
        key: `${key + id}cost`,
        icon: <IconFile />,
      });
    }
    return loadNodes;
  };

  const buildShardNodes = (total: any, key: string, shards: any) => {
    let children: any = [];
    const shardNodes = {
      label: <Text strong>Shards(Async)</Text>,
      key: `${key}shards`,
      children: children,
    };
    for (let shardID of Object.keys(shards)) {
      const shardStats = shards[shardID];
      let nodeStats = {
        label: (
          <span>
            <Text strong link>
              {shardID}{" "}
            </Text>
            : Series Filtering: [ Cost:{" "}
            {renderCost(shardStats.seriesFilterCost, total)}, Num. Of Series:{" "}
            <Text link>{formatter(shardStats.numOfSeries, UnitEnum.None)}</Text>{" "}
            ]
          </span>
        ),
        key: key + shardID,
        icon: <IconFlowChartStroked />,
        children: [
          {
            label: (
              <span>
                <Text strong> Memory Filtering</Text>:{" "}
                {renderCost(shardStats.memFilterCost, total)}
              </span>
            ),
            key: `${key + shardID}memory filtering`,
          },
          {
            label: (
              <span>
                <Text strong>KV Store Filtering:</Text>{" "}
                {renderCost(shardStats.kvFilterCost, total)}
              </span>
            ),
            key: `${key + shardID}kv store filtering`,
          },
          {
            label: (
              <span>
                <Text strong>Grouping</Text>:{" "}
                {renderCost(shardStats.groupingCost, total)}
              </span>
            ),
            key: `${key + shardID}Grouping`,
          },
          {
            label: (
              <span>
                <Text strong>Group Build(Async)</Text>: [ Count:{" "}
                <Text link>
                  {formatter(shardStats.groupBuildStats.count, UnitEnum.None)}
                </Text>
                , Total Cost:{" "}
                {renderCost(shardStats.groupBuildStats.totalCost, total)}, Min:{" "}
                {renderCost(shardStats.groupBuildStats.min, total)}, Max:{" "}
                {renderCost(shardStats.groupBuildStats.min, total)} ]
              </span>
            ),
            key: `${key + shardID}group build`,
          },
        ],
      };
      children.push(nodeStats);
      if (shardStats.scanStats) {
        nodeStats.children.push(
          buildLoadStats(total, key + shardID, shardStats.scanStats)
        );
      }
    }
    return shardNodes;
  };
  const buildStorageNodes = (total: any, storageNodes: any) => {
    let children: any = [];
    let storageNode = {
      label: <Text strong>Storage Nodes</Text>,
      key: "Storage-Nodes",
      children: children,
    };
    for (let key of Object.keys(storageNodes)) {
      const storageNodeStats = storageNodes[key];

      let nodeStats = {
        label: (
          <span>
            <Text strong link>
              {key}
            </Text>
            : [ Cost: {renderCost(storageNodeStats.totalCost, total)}, Network
            Payload:{" "}
            <Text link>
              {" "}
              {formatter(storageNodeStats.netPayload, UnitEnum.Bytes)})
            </Text>{" "}
            ]
          </span>
        ),
        key: key,
        icon: <IconServerStroked />,
        children: [
          {
            label: (
              <span>
                <Text strong>Execute Plan</Text>:{" "}
                {renderCost(storageNodeStats.planCost, total)}
              </span>
            ),
            key: `${key}plan-execute`,
          },
          {
            label: (
              <span>
                <Text strong>Tag Metadata Filtering</Text>:{" "}
                {renderCost(storageNodeStats.tagFilterCost, total)}
              </span>
            ),
            key: `${key}tag-filtering`,
          },
        ],
      };
      children.push(nodeStats);
      if (storageNodeStats.collectTagValuesStats) {
        nodeStats.children.push(
          buildCollectValueStats(
            total,
            key,
            storageNodeStats.collectTagValuesStats
          )
        );
      }
      if (storageNodeStats.shards) {
        nodeStats.children.push(
          buildShardNodes(total, key, storageNodeStats.shards)
        );
      }
    }
    return storageNode;
  };

  const buildStatsData = () => {
    if (!state) {
      return [];
    }
    let root = {
      label: (
        <>
          <Text strong>Root</Text>:{" "}
          <Text link>{formatter(state.totalCost, UnitEnum.Nanoseconds)}</Text>
        </>
      ),
      key: "Root",
      children: [
        {
          label: (
            <span>
              <Text strong>Execute Plan</Text> :{" "}
              {renderCost(state.planCost, state.totalCost)}
            </span>
          ),
          key: "Execute Plan",
        },
        {
          label: (
            <span>
              <Text strong>Waiting Storeage Response</Text>:{" "}
              {renderCost(state.waitCost, state.totalCost)}
            </span>
          ),
          key: "Waiting Storeage Response",
        },
        {
          label: (
            <span>
              <Text strong>Expression Eval</Text>:{" "}
              {renderCost(state.expressCost, state.totalCost)}
            </span>
          ),
          key: "Expression Eval",
        },
      ],
    };
    if (state.storageNodes) {
      root.children.push(
        buildStorageNodes(state.totalCost, state.storageNodes)
      );
    }
    return [root];
  };

  useEffect(() => {
    const disposer = [
      reaction(
        () => ChartStore.chartStatusMap.get(chartId),
        (s: ChartStatus | undefined) => {
          if (!s || s == ChartStatus.Loading) {
            return;
          }
          const state = ChartStore.stateCache.get(chartId);
          setState(state);
          console.log("explain state", state);
        }
      ),
    ];

    return () => {
      disposer.forEach((d) => d());
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  console.log("buildStatsData()", buildStatsData());

  return (
    <>
      <Tree
        expandAll
        // icon
        defaultExpandAll
        treeData={buildStatsData()}
        style={{ display: state ? "block" : "none" }}
      />
      <div style={{ display: !state ? "block" : "none" }}>
        <CanvasChart chartId={chartId} height={300} />
      </div>
    </>
  );
};

export default ExplainStatsView;
