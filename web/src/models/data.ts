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
export type Metadata = {
  type: "database" | "namespace" | "metric" | "field" | "tagKey" | "tagValue";
  values: string[] | Object[];
};

export type ResultSet = {
  metricName: string;
  startTime: number;
  endTime: number;
  interval: number;
  series: Series[];
  stats?: ExplainResult;
};
export type Series = {
  tags?: { [propName: string]: string };
  fields?: { [propName: string]: { [timestamp: number]: number } };
};

export type ExplainResult = {
  root: string;
  totalCost: number;
  netPayload: number;
  planCost: number;
  expressCost: number;
  waitCost: number;
  start: number;
  end: number;
  leafNodes: { [propName: string]: LeafNodeStats };
  brokerNodes: { [propName: string]: BrokerNodeExecStats };
};

export type BrokerNodeExecStats = {
  waitCost: number;
  storageNodes: { [propName: string]: LeafNodeStats };
};

export type LeafNodeStats = {
  netPayload: number;
  totalCost: number;
  start: number;
  end: number;
  stages: StageStats[];
};

export type StageStats = {
  identifier: string;
  start: number;
  end: number;
  async: boolean;
  cost: number;
  stage: string;
  errMsg: string;
  operators: OperatorStats[];
  children: StageStats[];
};
export type OperatorStats = {
  identifier: string;
  start: number;
  end: number;
  cost: number;
  errMsg: string;
  stats: any;
};

export type ShardExecStats = {
  groupingCost: string;
  kvFilterCost: string;
  memFilterCost: string;
  numOfSeries: number;
  seriesFilterCost: string;
  groupBuildStats: GroupBuildStats;
  scanStats: { [propName: string]: ScanStats };
};

export type GroupBuildStats = {
  count: number;
  min: string;
  max: string;
};

export type ScanStats = {
  count: number;
  min: string;
  max: string;
};
