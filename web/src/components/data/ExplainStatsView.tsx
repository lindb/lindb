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
  IconFixedStroked,
  IconShareStroked,
  IconServerStroked,
  IconTemplateStroked,
} from "@douyinfe/semi-icons";
import { Tree, Typography } from "@douyinfe/semi-ui";
import { CanvasChart, MetricStatus } from "@src/components";
import {
  ChartStatus,
  ExplainResult,
  LeafNodeStats,
  OperatorStats,
  StageStats,
  UnitEnum,
} from "@src/models";
import { ChartStore } from "@src/stores";
import { formatter } from "@src/utils";
import { reaction } from "mobx";
import React, { useEffect, useState } from "react";
import * as _ from "lodash-es";
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

  const buildOperatorStats = (
    parent: string,
    total: number,
    operators: OperatorStats[]
  ): any[] => {
    const children: any[] = [];
    _.forEach(operators, (item: OperatorStats, idx: number) => {
      const key = `${`${parent}-Leaf-Stage-Operator-${item.identifier}-${idx}`}`;
      const operatorNode = {
        label: (
          <span>
            <Text strong>{item.identifier}</Text>: [ Cost:{" "}
            {renderCost(item.cost, total)} ]
          </span>
        ),
        key: key,
        icon: <IconFixedStroked />,
      };
      children.push(operatorNode);
    });
    return children;
  };

  const buildStageStats = (
    parent: string,
    total: number,
    stages: StageStats[]
  ): any[] => {
    if (_.isEmpty(stages)) {
      return [];
    }
    let children: any = [];
    _.forEach(stages, (item: StageStats, idx: number) => {
      const key = `${`${parent}-Leaf-Stage-${item.identifier}-${idx}`}`;
      const stageNode = {
        label: (
          <span>
            <Text strong>{item.identifier}</Text>: [{" "}
            {item.async && <Text type="success">Async</Text>} Cost:{" "}
            {renderCost(item.cost, total)} ]
          </span>
        ),
        key: key,
        icon: <IconShareStroked />,
        children: [] as any[],
      };
      children.push(stageNode);
      if (!_.isEmpty(item.operators)) {
        stageNode.children.push(
          ...buildOperatorStats(key, total, item.operators)
        );
      }

      if (!_.isEmpty(item.children)) {
        stageNode.children.push(...buildStageStats(key, total, item.children));
      }
    });
    return children;
  };

  const buildLeafNodes = (parent: string, total: number, leafNodes: any) => {
    let children: any = [];
    for (let key of Object.keys(leafNodes)) {
      const leafNodeStats = leafNodes[key] as LeafNodeStats;
      const nKey = `${parent}-Leaf-Nodes-${key}`;

      let nodeStats = {
        label: (
          <span>
            <Text strong>
              Leaf[
              <Text strong link>
                {key}
              </Text>
              ]
            </Text>
            : [ Cost: {renderCost(leafNodeStats.totalCost, total)}, Network
            Payload:{" "}
            <Text link>
              {" "}
              {formatter(leafNodeStats.netPayload, UnitEnum.Bytes)}
            </Text>{" "}
            ]
          </span>
        ),
        key: nKey,
        icon: <IconServerStroked />,
        children: buildStageStats(nKey, total, leafNodeStats.stages),
      };
      children.push(nodeStats);
    }
    return children;
  };

  const buildIntermediateBrokers = (total: any, brokerNodes: any) => {
    let children: any = [];
    let IntermediateNode = {
      label: <Text strong>Intermediate Nodes</Text>,
      key: "Intermediate-Nodes",
      children: children,
    };
    for (let key of Object.keys(brokerNodes)) {
      const brokerNodeStats = brokerNodes[key];
      const nodeStats = {
        label: (
          <span>
            <Text strong link>
              {key}
            </Text>
            : [ Waiting: {renderCost(brokerNodeStats.waitCost, total)}, Cost:{" "}
            {renderCost(brokerNodeStats.totalCost, total)}, Network Payload:{" "}
            <Text link>
              {" "}
              {formatter(brokerNodeStats.netPayload, UnitEnum.Bytes)})
            </Text>{" "}
            ]
          </span>
        ),
        icon: <IconTemplateStroked />,
        key: `Intermediate-${key}`,
        children: [] as any,
      };
      children.push(nodeStats);
      if (brokerNodeStats.storageNodes) {
        nodeStats.children.push(
          buildLeafNodes(key, total, brokerNodeStats.storageNodes)
        );
      }
    }
    return IntermediateNode;
  };

  const buildStatsData = () => {
    if (!state) {
      return [];
    }
    let root = {
      label: (
        <>
          <Text strong>
            Root[
            <Text strong link>
              {state.root}
            </Text>
            ]
          </Text>
          : [ Cost:{" "}
          <Text link>{formatter(state.totalCost, UnitEnum.Nanoseconds)}</Text>,
          Network Payload:{" "}
          <Text link>{formatter(state.netPayload, UnitEnum.Bytes)}</Text> ]
        </>
      ),
      key: "Root",
      icon: <IconServerStroked />,
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
              <Text strong>Waiting Response</Text>:{" "}
              {renderCost(state.waitCost, state.totalCost)}
            </span>
          ),
          key: "Waiting Intermediate Response",
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
    if (state.brokerNodes) {
      root.children.push(
        buildIntermediateBrokers(state.totalCost, state.brokerNodes)
      );
    }
    if (state.leafNodes) {
      root.children.push(
        ...buildLeafNodes("root", state.totalCost, state.leafNodes)
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
        }
      ),
    ];

    return () => {
      disposer.forEach((d) => d());
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <MetricStatus chartId={chartId} showMsg />
      <Tree
        expandAll
        // icon
        defaultExpandAll
        treeData={buildStatsData()}
        style={{ display: state ? "block" : "none" }}
      />
      <div style={{ display: !state ? "block" : "none" }}>
        <CanvasChart chartId={chartId} height={300} disableDrag />
      </div>
    </>
  );
};

export default ExplainStatsView;
