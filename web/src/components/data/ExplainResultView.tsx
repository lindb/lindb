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
  StorageNodeExecStats,
  ShardExecStats,
  ScanStats,
} from "@src/models";
import { Tree, Typography } from "@douyinfe/semi-ui";
import { IconTemplate } from "@douyinfe/semi-icons";
const Text = Typography.Text;

interface ExplainResultViewProps {
  stats: ExplainResult;
}

export default function ExplainResultView(props: ExplainResultViewProps) {
  const { stats } = props;
  const renderLabel: React.FC<any> = ({
    className,
    onExpand,
    onClick,
    data,
    expandIcon,
  }) => {
    const { label, icon } = data;
    const isLeaf = !(data.children && data.children.length);
    return (
      <li
        className={className}
        role="treenode"
        onClick={isLeaf ? onClick : onExpand}
      >
        {/* {isLeaf ? null : expandIcon} */}
        {icon ? icon : null}
        <span>{label}</span>
      </li>
    );
  };

  const generateScan = (parent: any, key: string, value: ScanStats) => {
    console.log(value);
    parent.children.push({
      label: (
        <>
          {key} =&gt; Count: <Text link>{value.count}</Text>, Min:
          <Text link>{value.min}</Text>, Max:
          <Text link>{value.max}</Text>
        </>
      ),
      value: `${key}groupBuild`,
      key: `${parent.key}-${key}groupBuild`,
    });
  };

  const generateShard = (parent: any, key: string, value: ShardExecStats) => {
    const shard = {
      label: (
        <>
          Shard {key} =&gt; Memory Filter:
          <Text link>{value.memFilterCost}</Text>, File Filter:
          <Text link>{value.memFilterCost}</Text>, Series Filter:
          <Text link>{value.seriesFilterCost}</Text>, Num. Of Series:
          <Text link>{value.numOfSeries}</Text>, Group By:
          <Text link>{value.groupingCost}</Text>
        </>
      ),
      value: key,
      key: key,
      children: [] as any[],
    };
    parent.children.push(shard);
    if (value.groupBuildStats) {
      shard.children.push({
        label: (
          <>
            Group Build =&gt; Count:
            <Text link>{value.groupBuildStats.count}</Text>, Min:
            <Text link>{value.groupBuildStats.min}</Text>, Max:
            <Text link>{value.groupBuildStats.max}</Text>
          </>
        ),
        value: `${key}groupBuild`,
        key: `${key}groupBuild`,
      });
    }
    if (value.scanStats) {
      const scan = {
        label: "Data Scan",
        value: `${key}-scan`,
        key: `${key}-scan`,
        children: [],
      };
      shard.children.push(scan);
      const scanStats = value.scanStats;
      Object.keys(scanStats).forEach((key) => {
        generateScan(scan, key, scanStats[key]);
      });
    }
  };
  const generateStorageNode = (
    parent: any,
    key: string,
    value: StorageNodeExecStats
  ) => {
    const stoargeNode = {
      label: (
        <>
          {key} =&gt; Total:<Text link>{value.totalCost}</Text>, SQL Plan:
          <Text link>{value.planCost}</Text>, Tag Filter:
          <Text link>{value.tagFilterCost}</Text>, Net:
          <Text link>{value.netPayload}</Text>
        </>
      ),
      value: key,
      key: key,
      children: [],
    };
    parent.children.push(stoargeNode);
    const shards = value.shards || {};
    Object.keys(shards).forEach((key) => {
      generateShard(stoargeNode, key, shards[key]);
    });
  };
  const generateTree = () => {
    const storageList = {
      label: "Storage Nodes",
      value: "StorageList",
      key: "StorageList",
      icon: <IconTemplate />,
      children: [] as any[],
    };
    const root = {
      label: (
        <>
          Root Node(Total:<Text link>{stats.totalCost}</Text>, SQL Plan:
          <Text link>{stats.planCost}</Text>, Expression Eval:
          <Text link>{stats.expressCost}</Text>, Wait:{" "}
          <Text link>{stats.expressCost}</Text>)
        </>
      ),
      value: "Root",
      key: "Root",
      children: [storageList],
    };
    const storageNodes = stats.storageNodes || {};
    Object.keys(storageNodes).forEach((key) => {
      generateStorageNode(storageList, key, storageNodes[key]);
    });
    return [root];
  };

  return (
    <Tree expandAll treeData={generateTree()} renderFullLabel={renderLabel} />
  );
}
