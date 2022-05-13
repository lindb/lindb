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
  IconHelpCircleStroked,
  IconSourceControl,
  IconSearch,
  IconTrueFalseStroked,
  IconTick,
} from "@douyinfe/semi-icons";
import {
  Button,
  Card,
  Col,
  Modal,
  Row,
  SplitButtonGroup,
  Tooltip,
  Tree,
  List,
  Input,
  Spin,
  Typography,
} from "@douyinfe/semi-ui";
import { StorageStatusView } from "@src/components";
import { StateRoleName } from "@src/constants";
import { Storage } from "@src/models";
import { exec } from "@src/services";
import * as _ from "lodash-es";
import * as monaco from "monaco-editor";
import React, {
  MutableRefObject,
  ReactNode,
  useEffect,
  useRef,
  useState,
} from "react";
const { Text, Title } = Typography;

export type CompareViewProps = {
  source: any;
  nodes: any;
};

export function CompareView(props: CompareViewProps) {
  const { source, nodes } = props;
  const diffEditor = useRef() as MutableRefObject<any>;
  const diffEditorRef = useRef() as MutableRefObject<HTMLDivElement>;
  const [filter, setFilter] = useState<string>("");

  useEffect(() => {
    if (diffEditorRef.current && !diffEditor.current) {
      console.log("xxxxxx");
      // editor no init, create it
      diffEditor.current = monaco.editor.createDiffEditor(
        diffEditorRef.current,
        {
          theme: "lindb",
          readOnly: true,
        }
      );
    }
    var originalModel = monaco.editor.createModel(
      _.get(source, "data", "no data"),
      "json"
    );
    var modifiedModel = monaco.editor.createModel(
      _.get(nodes, "[0]data", "no data"),
      "json"
    );
    console.log("diffEditor.current", diffEditor.current);
    diffEditor.current.setModel({
      original: originalModel,
      modified: modifiedModel,
    });
  }, [source, nodes]);

  return (
    <>
      <Row gutter={12}>
        <Col span={4}>
          <List
            bordered
            header={
              <Input
                onChange={(v) => (!v ? setFilter(v) : null)}
                placeholder="Filter nodes"
                prefix={<IconSearch />}
              />
            }
            size="small"
            dataSource={_.filter(nodes, function (o) {
              console.log(
                "o.node.indexOf(filter)",
                filter,
                o.node.indexOf(filter)
              );
              return o.node.indexOf(filter) >= 0;
            })}
            renderItem={(item) => (
              <List.Item
                main={item.node}
                extra={
                  item.isDiff ? (
                    <IconTrueFalseStroked
                      style={{ color: "var(--semi-color-danger)" }}
                    />
                  ) : (
                    <IconTick style={{ color: "var(--semi-color-success)" }} />
                  )
                }
              />
            )}
          />
        </Col>
        <Col span={20}>
          <div ref={diffEditorRef} style={{ height: "90vh" }} />
        </Col>
      </Row>
    </>
  );
}

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
  const editor = useRef() as MutableRefObject<any>;
  const editorRef = useRef() as MutableRefObject<HTMLDivElement>;
  const [root, setRoot] = useState<any[]>([]);
  const [loadedKeys, setLoadedKeys] = useState<any[]>([]);
  const [metadata, setMetadata] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [comparing, setComparing] = useState(false);
  const [currentNode, setCurrentNode] = useState<any>(null);
  const [stateMachineMetadata, setStateMachineMetadata] = useState<any[]>([]);
  const [showCompareResult, setShowCompareResult] = useState<boolean>(false);

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
    if (editorRef.current && !editor.current) {
      // editor no init, create it
      editor.current = monaco.editor.create(editorRef.current, {
        value: "no data",
        language: "json",
        lineNumbers: "off",
        minimap: { enabled: false },
        lineNumbersMinChars: 2,
        readOnly: true,
        theme: "lindb",
      });
    }
    editor.current.setValue(_.get(metadata, "data", "no data"));
  }, [metadata]);

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
      setCurrentNode(node);
      setStateMachineMetadata([]);
      setLoading(true);
      let storageClause = "";
      if (node.storage) {
        storageClause = ` and storage='${node.storage}'`;
      }
      const metadata = await exec<any>({
        sql: `show ${node.role} metadata from state_repo where type='${node.type}'${storageClause}`,
      });
      setMetadata({
        type: node.type,
        data: JSON.stringify(metadata, null, "\t"),
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
    height: "83vh",
    border: "1px solid var(--semi-color-border)",
  };

  const exploreStateMachineData = async () => {
    try {
      const node = currentNode as Node;
      setComparing(true);
      setStateMachineMetadata([]);
      let storageClause = "";
      if (node.storage) {
        storageClause = ` and storage='${node.storage}'`;
      }
      const metadataFromSM = await exec<any>({
        sql: `show ${node.role} metadata from state_machine where type='${node.type}'${storageClause}`,
      });
      const nodes: any[] = [];
      _.mapKeys(metadataFromSM, function (data, key) {
        const dataFromSM = JSON.stringify(data, null, "\t");
        nodes.push({
          node: key,
          data: dataFromSM,
          isDiff: !_.isEqual(dataFromSM, metadata.data),
        });
      });
      setStateMachineMetadata(nodes);
    } finally {
      setComparing(false);
    }
  };

  return (
    <>
      <Row gutter={8}>
        <Col span={8}>
          <Card>
            <Tree
              loadedKeys={loadedKeys}
              treeData={root}
              renderFullLabel={renderLabel}
              style={treeStyle}
              onChange={(args: any) => loadMetadata(args as Node)}
            />
          </Card>
        </Col>
        <Col span={16} style={{ display: "flex", flexDirection: "column" }}>
          <div>
            <SplitButtonGroup style={{ marginBottom: 8, marginRight: 12 }}>
              <Button
                icon={<IconSourceControl />}
                onClick={exploreStateMachineData}
              >
                Compare
              </Button>
              <Tooltip content="Compare with state matchine's data in memory">
                <Button icon={<IconHelpCircleStroked />} />
              </Tooltip>
            </SplitButtonGroup>
            {comparing && (
              <>
                <Spin size="middle" />
                <Text style={{ marginRight: 4 }}>Comparing</Text>
              </>
            )}
            {!_.isEmpty(stateMachineMetadata) && (
              <Text strong link onClick={() => setShowCompareResult(true)}>
                Found <Text type="success">{stateMachineMetadata.length}</Text>{" "}
                nodes, diff{" "}
                <Text type="danger">
                  {_.filter(stateMachineMetadata, (o) => o.isDiff).length}
                </Text>{" "}
                nodes.
              </Text>
            )}
          </div>

          <Card style={{ height: "83.5vh" }}>
            <div ref={editorRef} style={{ height: "80vh" }} />
          </Card>
        </Col>
      </Row>

      <Modal
        title={
          <div>
            <Title strong heading={4}>
              State Compare Result
            </Title>
            <Text>
              The state(in storage) compare with in state machine(memory)
            </Text>
          </div>
        }
        closeOnEsc
        fullScreen
        footer={null}
        visible={showCompareResult}
        onCancel={() => {
          setShowCompareResult(false);
        }}
      >
        <CompareView source={metadata} nodes={stateMachineMetadata} />
      </Modal>
    </>
  );
}
