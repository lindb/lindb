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
import React, { useEffect, useState } from "react";
import { Card, Row, Col, Tree } from "@douyinfe/semi-ui";
import { metaExplore, exploreRepoData } from "@src/services";
import { StateRoleName } from "@src/constants";
import * as _ from "lodash-es";

export default function Explore() {
  const [root, setRoot] = useState<any[]>([]);
  const [loadedKeys, setLoadedKeys] = useState<any[]>([]);
  const getItems = (parent: string, obj: any) => {
    const keys = _.keys(obj);
    const rs: any[] = [];
    _.forEach(keys, (k) =>
      rs.push({
        label: k,
        value: k,
        key: `${parent}-${k}`,
        parent: parent,
        data: _.get(obj, k, {}),
      })
    );
    return rs;
  };
  useEffect(() => {
    const fetchMetadata = async () => {
      const metadata = await metaExplore();
      const keys = _.keys(metadata);
      const root: any[] = [];
      const loadedKeys: any[] = [];
      _.forEach(keys, (key) => {
        if (key !== StateRoleName.Storage) {
          const data = _.get(metadata, key, {});
          loadedKeys.push(key);
          root.push({
            label: key,
            value: key,
            key: key,
            data: data,
            children: getItems(key, data),
          });
        }
      });
      setRoot(root);
    };
    fetchMetadata();
  }, []);

  const loadMetadata = async (node: any) => {
    console.log("load...", node);
    if (node.children) {
      return;
    }
    await exploreRepoData({ role: node.parent, type: node.value });
  };

  const renderLabel: React.FC<any> = ({
    className,
    onExpand,
    onClick,
    data,
    expandIcon,
  }) => {
    const { label } = data;
    const isLeaf = !(data.children && data.children.length);
    return (
      <li
        className={className}
        role="treenode"
        onClick={isLeaf ? onClick : onExpand}
      >
        {isLeaf ? null : expandIcon}
        <span>{label}</span>
      </li>
    );
  };
  const treeStyle = {
    width: "100%",
    height: "80vh",
    border: "1px solid var(--semi-color-border)",
  };

  return (
    <>
      <Row gutter={8}>
        <Col span={8}>
          <Card bordered={true}>
            <Tree
              //   directory
              // blockNode
              // expandedKeys={loadedKeys}
              loadedKeys={loadedKeys}
              loadData={loadMetadata}
              treeData={root}
              // renderFullLabel={renderLabel}
              style={treeStyle}
              onChange={(...args) => console.log("change", ...args)}
            />
          </Card>
        </Col>
        <Col span={16}>
          <Card bordered={false}>content</Card>
        </Col>
      </Row>
    </>
  );
}
