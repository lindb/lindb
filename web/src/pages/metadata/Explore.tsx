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
import { IconHelpCircleStroked, IconSourceControl } from "@douyinfe/semi-icons";
import {
  Button,
  Card,
  Col,
  Row,
  SplitButtonGroup,
  Tooltip,
  Tree,
} from "@douyinfe/semi-ui";
import Editor from "@monaco-editor/react";
import { StorageStatusView } from "@src/components";
import { StateRoleName } from "@src/constants";
import { Storage } from "@src/models";
import { exec } from "@src/services";
import * as _ from "lodash-es";
import React, { ReactNode, useEffect, useState } from "react";

type Node = {
  role: string;
  type: string;
  storage?: string;
};

type TreeNode = {
  label: string | ReactNode;
  value: string;
  key: string;
  parent: StateRoleName;
  children: any[];
};

export default function MetadataExplore() {
  const [root, setRoot] = useState<any[]>([]);
  const [loadedKeys, setLoadedKeys] = useState<any[]>([]);
  const [metadata, setMetadata] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const getItems = (
    parent: string,
    role: string,
    obj: any,
    storage?: string
  ) => {
    const keys = _.keys(obj);
    const rs: any[] = [];
    _.forEach(keys, (k) =>
      rs.push({
        label: k,
        value: { role: role, type: k, storage: storage },
        key: `${parent}-${k}`,
        parent: parent,
      })
    );
    return rs;
  };
  useEffect(() => {
    const fetchMetadata = async () => {
      const metadata = await exec<any>({ sql: "show metadata types" });
      const storages = await exec<Storage[]>({ sql: "show storages" });
      const keys = _.keys(metadata);
      _.sortBy(keys);
      const root: any[] = [];
      const loadedKeys: any[] = [];
      _.forEach(keys, (key) => {
        const data = _.get(metadata, key, {});
        loadedKeys.push(key);
        if (key === StateRoleName.Storage) {
          const storageNodes: TreeNode[] = [];
          _.forEach(storages || [], (storage: any) => {
            const namespace = storage.config.namespace;
            const storageKey = `${key}-${namespace}`;
            const storageTypes = getItems(storageKey, key, data, namespace);
            storageNodes.push({
              label: (
                <>
                  {namespace} (<StorageStatusView text={storage.status} />)
                </>
              ),
              value: namespace,
              key: storageKey,
              parent: key,
              children: storageTypes,
            });
          });
          root.push({
            label: key,
            value: key,
            key: key,
            children: storageNodes,
          });
        } else {
          root.push({
            label: key,
            value: key,
            key: key,
            children: getItems(key, key, data),
          });
        }
      });
      setLoadedKeys(loadedKeys);
      setRoot(root);
    };
    fetchMetadata();
  }, []);

  const loadMetadata = async (node: Node) => {
    try {
      setLoading(true);
      let storageClause = "";
      if (node.storage) {
        storageClause = ` and storage='${node.storage}'`;
      }
      const metadata = await exec<any>({
        sql: `show ${node.role} metadata where type='${node.type}'${storageClause}`,
      });
      setMetadata({
        type: node.type,
        data: metadata,
        node: node,
      });
    } finally {
      setLoading(false);
    }
  };

  const renderLabel: React.FC<any> = ({
    className,
    onExpand,
    onClick,
    data,
    expandIcon,
  }) => {
    const { label } = data;
    let isLeaf = !(data.children && data.children.length);
    return (
      <li
        className={className}
        role="treeitem"
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
              loadedKeys={loadedKeys}
              treeData={root}
              renderFullLabel={renderLabel}
              style={treeStyle}
              onChange={(args: any) => loadMetadata(args as Node)}
            />
          </Card>
        </Col>
        <Col span={16}>
          <SplitButtonGroup style={{ marginBottom: 8 }}>
            <Button icon={<IconSourceControl />}>Compare</Button>
            <Tooltip content="Compare with state matchine">
              <Button icon={<IconHelpCircleStroked />} />
            </Tooltip>
          </SplitButtonGroup>
          <Card style={treeStyle} bordered={false}>
            <Editor
              theme="vs-dark"
              height="90vh"
              defaultLanguage="json"
              options={{
                readOnly: true,
                minimap: { enabled: false },
                lineNumbersMinChars: 2,
              }}
              value={JSON.stringify(
                _.get(metadata, "data", "no data"),
                null,
                "\t"
              )}
              loading={loading}
            />
          </Card>
        </Col>
      </Row>
    </>
  );
}
