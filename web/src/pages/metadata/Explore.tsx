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
  Empty,
} from "@douyinfe/semi-ui";
import { Icon, StatusTip, StorageStatusView } from "@src/components";
import { StateRoleName, Theme } from "@src/constants";
import { UIContext } from "@src/context/UIContextProvider";
import { Storage } from "@src/models";
import { ExecService } from "@src/services";
import { useQuery } from "@tanstack/react-query";
import * as _ from "lodash-es";
import * as monaco from "monaco-editor";
import React, {
  MutableRefObject,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
const { Text, Title } = Typography;

const CompareView: React.FC<{
  source: any;
  nodes: any;
}> = (props) => {
  const { source, nodes } = props;
  const diffEditor = useRef() as MutableRefObject<any>;
  const diffEditorRef = useRef() as MutableRefObject<HTMLDivElement>;
  const [filter, setFilter] = useState<string>("");
  const [current, setCurrent] = useState<string>("");
  const { theme } = useContext(UIContext);

  const buildOriginal = useCallback(() => {
    return monaco.editor.createModel(_.get(source, "data", "no data"), "json");
  }, [source]);

  useEffect(() => {
    if (diffEditorRef.current && !diffEditor.current) {
      // editor no init, create it
      diffEditor.current = monaco.editor.createDiffEditor(
        diffEditorRef.current,
        {
          theme: theme === Theme.dark ? "vs-dark" : "vs",
          readOnly: true,
        }
      );
    }
    var modifiedModel = monaco.editor.createModel(
      _.get(nodes, "[0].data", "no data"),
      "json"
    );
    diffEditor.current.setModel({
      original: buildOriginal(),
      modified: modifiedModel,
    });
    setCurrent(_.get(nodes, "[0].node", "-"));
  }, [source, nodes, buildOriginal, theme]);

  return (
    <>
      <Row gutter={12}>
        <Col span={4}>
          <List
            bordered
            header={
              <Input
                onChange={setFilter}
                placeholder="Filter nodes"
                prefix={<IconSearch />}
              />
            }
            size="small"
            dataSource={_.filter(nodes, function (o) {
              return o.node.indexOf(filter) >= 0;
            })}
            renderItem={(item) => (
              <List.Item
                style={{
                  cursor: "pointer",
                  backgroundColor:
                    item.node === current
                      ? "var(--semi-color-primary-light-default)"
                      : "",
                }}
                main={item.node}
                className="list-item"
                onClick={() => {
                  setCurrent(item.node);
                  var modifiedModel = monaco.editor.createModel(
                    item.data,
                    "json"
                  );
                  diffEditor.current.setModel({
                    modified: modifiedModel,
                    original: buildOriginal(),
                  });
                }}
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
};

const MetadataView: React.FC<{
  node: Node;
}> = (props) => {
  const { node } = props;
  const editorRef = useRef() as MutableRefObject<HTMLDivElement>;
  const [comparing, setComparing] = useState(false);
  const [stateMachineMetadata, setStateMachineMetadata] = useState<any[]>([]);
  const [showCompareResult, setShowCompareResult] = useState<boolean>(false);
  const [metadata, setMetadata] = useState<any>(null);
  const { theme } = useContext(UIContext);

  const { isLoading, isInitialLoading, isError, error, data } = useQuery(
    ["load_metadata", node],
    async () => {
      setStateMachineMetadata([]);
      let storageClause = "";
      if (node.storage) {
        storageClause = ` and storage='${node.storage}'`;
      }
      return ExecService.exec<any>({
        sql: `show ${node.role} metadata from state_repo where type='${node.type}'${storageClause}`,
      });
    },
    {
      enabled: node != null,
    }
  );

  useEffect(() => {
    const metadata = {
      type: node?.type,
      data: JSON.stringify(data, null, "\t"),
      node: node,
    };
    if (editorRef.current) {
      monaco.editor.create(editorRef.current, {
        theme: theme === Theme.dark ? "vs-dark" : "vs",
        language: "json",
        lineNumbers: "off",
        minimap: { enabled: false },
        lineNumbersMinChars: 2,
        readOnly: true,
        value: _.get(metadata, "data", "no data"),
      });
    }
    setMetadata(metadata);
  }, [data, node, theme]);

  const exploreStateMachineData = async () => {
    try {
      setComparing(true);
      setStateMachineMetadata([]);
      let storageClause = "";
      if (node.storage) {
        storageClause = ` and storage='${node.storage}'`;
      }
      const metadataFromSM = await ExecService.exec<any>({
        sql: `show ${node.role} metadata from state_machine where type='${node.type}'${storageClause}`,
      });
      const nodes: any[] = [];
      _.mapKeys(metadataFromSM, function (val, key) {
        const dataFromSM = JSON.stringify(val, null, "\t");
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

  const renderData = () => {
    if (!node) {
      return (
        <Empty
          image={<Icon icon="iconempty" style={{ fontSize: 48 }} />}
          description="No data"
          style={{ paddingTop: 100 }}
        />
      );
    }
    if (isLoading || isInitialLoading || isError) {
      return (
        <StatusTip
          style={{ paddingTop: 100 }}
          isLoading={isLoading || isInitialLoading}
          isError={isError}
          error={error}
        />
      );
    }
    return (
      <>
        <div ref={editorRef} style={{ height: "100%" }} />
      </>
    );
  };

  return (
    <>
      {node && (
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
      )}

      <Card>
        <div style={{ height: "80vh" }}>{renderData()}</div>
      </Card>

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
};

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

const MetadataExplore: React.FC = () => {
  const [root, setRoot] = useState<any[]>([]);
  const [loadedKeys, setLoadedKeys] = useState<any[]>([]);
  const [currentNode, setCurrentNode] = useState<any>(null);

  const { isLoading, isError, error, data } = useQuery(
    ["show_metadata"],
    async () => {
      return Promise.allSettled([
        ExecService.exec<any>({
          sql: "show metadata types",
        }),
        ExecService.exec<Storage[]>({
          sql: "show storages",
        }),
      ]).then((res) => {
        return res.map((item) =>
          item.status === "fulfilled" ? item.value : []
        );
      });
    }
  );

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
    if ((data || []).length != 2) {
      return;
    }
    const metadata = data![0];
    const storages = data![1];
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
  }, [data]);

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

  const renderMetadata = () => {
    if (isLoading || isError) {
      return (
        <StatusTip
          style={{ marginTop: 50 }}
          isLoading={isLoading}
          isError={isError}
          error={error}
        />
      );
    }
    return (
      <Tree
        className="lin-tree"
        loadedKeys={loadedKeys}
        treeData={root}
        renderFullLabel={renderLabel}
        onChange={(args: any) => setCurrentNode(args as Node)}
      />
    );
  };

  return (
    <>
      <Row gutter={8} style={{ display: "flex" }}>
        <Col span={8}>
          <Card bodyStyle={{ padding: 12 }} style={{ height: "100%" }}>
            {renderMetadata()}
          </Card>
        </Col>
        <Col span={16} style={{ display: "flex", flexDirection: "column" }}>
          <MetadataView node={currentNode} />
        </Col>
      </Row>
    </>
  );
};

export default MetadataExplore;
