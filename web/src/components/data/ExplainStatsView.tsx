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
import React, { useEffect, useState, useRef, MutableRefObject } from "react";
import * as _ from "lodash-es";
const Text = Typography.Text;

interface ExplainStatsViewProps {
  chartId: string;
}

type totalStats = {
  total: number;
  start: number;
  end: number;
};

type stats = {
  start: number;
  end: Number;
  cost: number;
};

const ExplainStatsView: React.FC<ExplainStatsViewProps> = (
  props: ExplainStatsViewProps
) => {
  const { chartId } = props;
  const [state, setState] = useState<ExplainResult>();
  const totalStats = useRef() as MutableRefObject<totalStats>;

  const getColor = (percent: number) => {
    let type: any = "success";
    if (percent > 50) {
      type = "danger";
    } else if (percent > 30) {
      type = "warning";
    }
    return type;
  };

  const timeline = (stats: stats, color?: string) => {
    const totalS = totalStats.current;
    const percent = (stats.cost * 100) / totalS.total;
    const offset = ((stats.start - totalS.start) * 100) / totalS.total;
    const background = color ? color : getColor(percent);
    return (
      <>
        <div className="lin-explain-timeline">
          <div
            className="inner"
            style={{
              backgroundColor: `var(--semi-color-${background})`,
              width: `${percent}%`,
              marginLeft: `calc(${offset}%)`,
            }}
          ></div>
        </div>
      </>
    );
  };
  const style = {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  };

  const renderCost = (cost: any) => {
    const percent = (cost * 100.0) / totalStats.current.total;
    let type: any = getColor(percent);
    return <Text type={type}>{formatter(cost, UnitEnum.Nanoseconds)}</Text>;
  };

  const buildOperatorStats = (
    parent: string,
    operators: OperatorStats[]
  ): any[] => {
    const children: any[] = [];
    _.forEach(operators, (item: OperatorStats, idx: number) => {
      const key = `${`${parent}-Leaf-Stage-Operator-${item.identifier}-${idx}`}`;
      const operatorNode = {
        label: (
          <div style={style}>
            <span>
              <Text strong>{item.identifier}</Text>: [ Cost:{" "}
              {renderCost(item.cost)} ]
            </span>
            {timeline({
              start: item.start,
              end: item.end,
              cost: item.cost,
            })}
          </div>
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
          <div style={style}>
            <span>
              <Text strong>{item.identifier}</Text>: [{" "}
              {item.async && <Text type="success">Async</Text>} Cost:{" "}
              {renderCost(item.cost)} ]
            </span>
            {timeline({ start: item.start, end: item.end, cost: item.cost })}
          </div>
        ),
        key: key,
        icon: <IconShareStroked />,
        children: [] as any[],
      };
      children.push(stageNode);
      if (!_.isEmpty(item.operators)) {
        stageNode.children.push(...buildOperatorStats(key, item.operators));
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
          <div style={style}>
            <span>
              <Text strong>
                Leaf[
                <Text strong link>
                  {key}
                </Text>
                ]
              </Text>
              : [ Cost: {renderCost(leafNodeStats.totalCost)}, Network Payload:{" "}
              <Text link>
                {" "}
                {formatter(leafNodeStats.netPayload, UnitEnum.Bytes)}
              </Text>{" "}
              ]
            </span>
            {timeline({
              start: leafNodeStats.start,
              end: leafNodeStats.end,
              cost: leafNodeStats.totalCost,
            })}
          </div>
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
          <div style={style}>
            <span>
              <Text strong link>
                {key}
              </Text>
              : [ Waiting: {renderCost(brokerNodeStats.waitCost)}, Cost:{" "}
              {renderCost(brokerNodeStats.totalCost)}, Network Payload:{" "}
              <Text link>
                {" "}
                {formatter(brokerNodeStats.netPayload, UnitEnum.Bytes)})
              </Text>{" "}
              ]
            </span>
            {timeline}
          </div>
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
    totalStats.current = {
      total: state.totalCost,
      start: state.start,
      end: state.end,
    };
    let root = {
      label: (
        <div style={style}>
          <span>
            <Text strong>
              Root[
              <Text strong link>
                {state.root}
              </Text>
              ]
            </Text>
            : [ Cost:{" "}
            <Text link>{formatter(state.totalCost, UnitEnum.Nanoseconds)}</Text>
            , Network Payload:{" "}
            <Text link>{formatter(state.netPayload, UnitEnum.Bytes)}</Text> ]
          </span>
          {timeline(
            {
              start: state.start,
              end: state.end,
              cost: state.totalCost,
            },
            "success"
          )}
        </div>
      ),
      key: "Root",
      icon: <IconServerStroked />,
      children: [
        {
          label: (
            <div style={style}>
              <span>
                <Text strong>Execute Plan</Text> : {renderCost(state.planCost)}
              </span>
              {timeline({
                start: state.planStart,
                end: state.planEnd,
                cost: state.planCost,
              })}
            </div>
          ),
          key: "Execute Plan",
        },
        {
          label: (
            <div style={style}>
              <span>
                <Text strong>Waiting Response</Text>:{" "}
                {renderCost(state.waitCost)}
              </span>
              {timeline({
                start: state.waitStart,
                end: state.waitEnd,
                cost: state.waitCost,
              })}
            </div>
          ),
          key: "Waiting Intermediate Response",
        },
        {
          label: (
            <div style={style}>
              <span>
                <Text strong>Expression Eval</Text>:{" "}
                {renderCost(state.expressCost)}
              </span>
              {timeline({
                start: state.expressStart,
                end: state.expressEnd,
                cost: state.expressCost,
              })}
            </div>
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
