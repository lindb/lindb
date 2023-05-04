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
} from "@douyinfe/semi-icons";
import { Tree, Typography } from "@douyinfe/semi-ui";
import { NodeStats, OperatorStats, StageStats, Unit } from "@src/models";
import { FormatKit } from "@src/utils";
import React, { useRef, MutableRefObject } from "react";
import * as _ from "lodash-es";
const Text = Typography.Text;

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

const ExplainStatsView: React.FC<{ state: NodeStats }> = (props) => {
  const { state } = props;
  const totalStats = useRef() as MutableRefObject<totalStats>;
  const nodeID = useRef(0);

  const getNodeID = (): string => {
    nodeID.current++;
    return `node-${nodeID.current}`;
  };

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
    return <Text type={type}>{FormatKit.format(cost, Unit.Nanoseconds)}</Text>;
  };

  const buildOperatorStats = (operators: OperatorStats[]): any[] => {
    const children: any[] = [];
    _.forEach(operators, (item: OperatorStats) => {
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
        key: getNodeID(),
        icon: <IconFixedStroked />,
      };
      children.push(operatorNode);
    });
    return children;
  };

  const buildStageStats = (stages: StageStats[]): any[] => {
    if (_.isEmpty(stages)) {
      return [];
    }
    let children: any = [];
    _.forEach(stages, (item: StageStats) => {
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
        key: getNodeID(),
        icon: <IconShareStroked />,
        children: [] as any[],
      };
      children.push(stageNode);
      if (!_.isEmpty(item.operators)) {
        stageNode.children.push(...buildOperatorStats(item.operators));
      }

      if (!_.isEmpty(item.children)) {
        stageNode.children.push(...buildStageStats(item.children));
      }
    });
    return children;
  };

  const buildNodeStats = (nodeStats: NodeStats): any => {
    const costs = [
      <span key="cost">
        <span style={{ marginRight: 2 }}>Cost:</span>
        <Text link>
          {FormatKit.format(nodeStats.totalCost, Unit.Nanoseconds)}
        </Text>
      </span>,
    ];
    if (nodeStats.waitStart > 0) {
      costs.push(
        <span key="wait">
          <span style={{ marginRight: 2 }}>, Wait:</span>
          <Text link>
            {FormatKit.format(nodeStats.waitCost, Unit.Nanoseconds)}
          </Text>
        </span>
      );
    }
    if (nodeStats.netPayload > 0) {
      costs.push(
        <span key="payload">
          <span style={{ marginRight: 2 }}>, Network:</span>
          <Text link>{FormatKit.format(nodeStats.netPayload, Unit.Bytes)}</Text>
        </span>
      );
    }
    let node = {
      label: (
        <div style={style}>
          <span>
            <Text strong>
              <Text strong link>
                {nodeStats.node}
              </Text>
            </Text>
            : [ {costs} ]
          </span>
          {timeline(
            {
              start: nodeStats.start,
              end: nodeStats.end,
              cost: nodeStats.totalCost,
            },
            nodeStats === state ? "success" : ""
          )}
        </div>
      ),
      key: getNodeID(),
      icon: <IconServerStroked />,
      children: [] as any[],
    };

    if (nodeStats.stages) {
      node.children.push(...buildStageStats(nodeStats.stages));
    }

    if (nodeStats.children) {
      nodeStats.children.forEach((child: NodeStats) => {
        node.children.push(buildNodeStats(child));
      });
    }
    return node;
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
    return [buildNodeStats(state)];
  };

  return (
    <>
      <Tree
        expandAll
        // icon
        defaultExpandAll
        treeData={buildStatsData()}
        style={{ display: state ? "block" : "none" }}
      />
    </>
  );
};

export default ExplainStatsView;
