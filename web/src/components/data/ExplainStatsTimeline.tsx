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
import {
  ExplainResult,
  LeafNodeStats,
  OperatorStats,
  StageStats,
} from "@src/models";
import { Progress } from "@douyinfe/semi-ui";
import * as _ from "lodash-es";

interface ExplainStatsTimelineProps {
  stats: ExplainResult | undefined;
}
const ExplainStatsTimeline: React.FC<ExplainStatsTimelineProps> = (
  props: ExplainStatsTimelineProps
) => {
  const { stats } = props;

  const buildOperators = (
    total: number,
    start: number,
    end: number,
    operators: OperatorStats[]
  ): any[] => {
    const operatorNodes: any[] = [];
    _.forEach(operators, (op: OperatorStats, idx: number) => {
      const nKey = `${idx}`;
      const percent = (op.cost * 100) / total;
      const p = ((op.start - start) * 100) / (end - start);
      console.log("operator...", p, percent);
      operatorNodes.push(
        <div
          key={nKey}
          className="lin-explain-timeline "
          style={{
            width: `{percent}%`,
            minWidth: 1,
            height: 12,
            backgroundColor: "red",
            // marginLeft: 10,
            marginLeft: `calc(${p}%)`,
          }}
        />
      );
    });
    return operatorNodes;
  };

  const buildStages = (
    total: number,
    start: number,
    end: number,
    stages: StageStats[]
  ): any[] => {
    const stageNodes: any[] = [];
    _.forEach(stages, (stage: StageStats, idx: number) => {
      const nKey = `${idx}`;
      const percent = (stage.cost * 100) / total;
      const p = ((stage.start - start) * 100) / (end - start);
      console.log("stage...", p, percent);
      stageNodes.push(
        <Progress
          key={nKey}
          className="lin-explain-timeline "
          style={{
            height: 12,
            backgroundColor: "rgba(255,255,255,0)",
            marginLeft: `calc(${p}%)`,
          }}
          percent={percent}
        />
      );
      if (!_.isEmpty(stage.operators)) {
        stageNodes.push(
          <div style={{ display: "flex", width: "100%" }}>
            {buildOperators(total, start, end, stage.operators)}
          </div>
        );
      }
      if (!_.isEmpty(stage.children)) {
        stageNodes.push(...buildStages(total, start, end, stage.children));
      }
    });
    return stageNodes;
  };
  const buildLeafNodes = (
    total: number,
    start: number,
    end: number,
    leafNodes: any
  ) => {
    const leaves = [];
    for (let key of Object.keys(leafNodes)) {
      const leafNodeStats = leafNodes[key] as LeafNodeStats;
      const nKey = `Leaf-Nodes-${key}`;
      const percent = (leafNodeStats.totalCost * 100) / total;
      const p = ((leafNodeStats.start - start) * 100) / (end - start);
      console.log("leaf...", p, percent);
      leaves.push(
        <Progress
          key={nKey}
          className="lin-explain-timeline "
          style={{
            height: 12,
            backgroundColor: "rgba(255,255,255,0)",
            marginLeft: `calc(${p}%)`,
          }}
          percent={percent}
        />
      );
      if (!_.isEmpty(leafNodeStats.stages)) {
        leaves.push(...buildStages(total, start, end, leafNodeStats.stages));
      }
    }
    return leaves;
  };
  const buildTimeline = () => {
    if (!stats) {
      return null;
    }
    const total = stats.totalCost;
    const start = stats.start;
    const end = stats.end;
    return (
      <>
        <Progress
          className="lin-explain-timeline "
          style={{ height: 12 }}
          percent={100}
        />
        {stats.leafNodes && buildLeafNodes(total, start, end, stats.leafNodes)}
      </>
    );
  };
  return <div>{buildTimeline()}</div>;
};

export default ExplainStatsTimeline;
