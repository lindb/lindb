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
import { Icon, StatusTip, ClusterStatusView } from "@src/components";
import { StateRoleName, Theme } from "@src/constants";
import { UIContext } from "@src/context/UIContextProvider";
import { Broker } from "@src/models";
import { ExecService } from "@src/services";
import { useQuery } from "@tanstack/react-query";
import {
  get,
  filter,
  mapKeys,
  isEqual,
  isEmpty,
  keys,
  forEach,
  sortBy,
} from "lodash-es";
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
  const [filterStr, setFilterStr] = useState<string>("");
  const [current, setCurrent] = useState<string>("");
  const { locale, theme } = useContext(UIContext);
  const { MetadataExploreView } = locale;

  const buildOriginal = useCallback(() => {
    return monaco.editor.createModel(get(source, "data", "no data"), "json");
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
      get(nodes, "[0].data", "no data"),
      "json"
    );
    diffEditor.current.setModel({
      original: buildOriginal(),
      modified: modifiedModel,
    });
    setCurrent(get(nodes, "[0].node", "-"));
  }, [source, nodes, buildOriginal, theme]);

  return (
    <>
      <Row gutter={12}>
        <Col span={4}>
          <List
            bordered
            header={
              <Input
                onChange={setFilterStr}
                placeholder={MetadataExploreView.filterNode}
                prefix={<IconSearch />}
              />
            }
            size="small"
            dataSource={filter(nodes, function (o) {
              return o.node.indexOf(filterStr) >= 0;
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
  const { theme, locale } = useContext(UIContext);
  const { Common, MetadataExploreView } = locale;

  const { isLoading, isInitialLoading, isError, error, data } = useQuery(
    ["load_metadata", node],
    async () => {
      setStateMachineMetadata([]);
      let storageClause = "";
      if (node.cluster) {
        if (node.role === StateRoleName.Broker) {
          return ExecService.exec<any>({
            sql: `show ${node.role} metadata from state_machine where type='${node.type}' and broker='${node.cluster}'`,
          });
        } else {
          storageClause = ` and storage='${node.cluster}'`;
        }
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
        value: get(metadata, "data", "no data"),
      });
    }
    setMetadata(metadata);
  }, [data, node, theme]);

  const exploreStateMachineData = async () => {
    try {
      setComparing(true);
      setStateMachineMetadata([]);
      let storageClause = "";
      if (node.cluster) {
        storageClause = ` and storage='${node.cluster}'`;
      }
      const metadataFromSM = await ExecService.exec<any>({
        sql: `show ${node.role} metadata from state_machine where type='${node.type}'${storageClause}`,
      });
      const nodes: any[] = [];
      mapKeys(metadataFromSM, function (val, key) {
        const dataFromSM = JSON.stringify(val, null, "\t");
        nodes.push({
          node: key,
          data: dataFromSM,
          isDiff: !isEqual(dataFromSM, metadata.data),
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
          description={Common.noData}
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
      {node &&
        (!node.cluster ||
          (node.cluster && node.role !== StateRoleName.Broker)) && (
          <div>
            <SplitButtonGroup style={{ marginBottom: 8, marginRight: 12 }}>
              <Button
                icon={<IconSourceControl />}
                onClick={exploreStateMachineData}
              >
                {MetadataExploreView.compare}
              </Button>
              <Tooltip content={MetadataExploreView.compareTooltip}>
                <Button icon={<IconHelpCircleStroked />} />
              </Tooltip>
            </SplitButtonGroup>
            {comparing && (
              <>
                <Spin size="middle" />
                <Text style={{ marginRight: 4 }}>
                  {MetadataExploreView.comparing}
                </Text>
              </>
            )}
            {!isEmpty(stateMachineMetadata) && (
              <Text strong link onClick={() => setShowCompareResult(true)}>
                {MetadataExploreView.compareResult1}{" "}
                <Text type="success">{stateMachineMetadata.length}</Text>{" "}
                {MetadataExploreView.compareResult2}{" "}
                <Text type="danger">
                  {filter(stateMachineMetadata, (o) => o.isDiff).length}
                </Text>{" "}
                {MetadataExploreView.compareResult3}
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
              {MetadataExploreView.compareResultTitle}
            </Title>
            <Text>{MetadataExploreView.compareResultDesc}</Text>
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
  cluster?: string;
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
  const { env } = useContext(UIContext);

  const { isLoading, isError, error, data } = useQuery(
    ["show_metadata"],
    async () => {
      const requests = [
        ExecService.exec<any>({
          sql: "show metadata types",
        }),
      ];
      if (env.role === StateRoleName.Root) {
        requests.push(
          ExecService.exec<Broker[]>({
            sql: "show brokers",
          })
        );
      }
      return Promise.allSettled(requests).then((res) => {
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
    cluster?: string
  ) => {
    const objKeys = keys(obj);
    const rs: any[] = [];
    forEach(objKeys, (k: any) =>
      rs.push({
        label: k,
        value: { role: role, type: k, cluster: cluster },
        key: `${parent}-${k}`,
        parent: parent,
      })
    );
    return rs;
  };

  useEffect(() => {
    if (!data) {
      return;
    }
    const root: any[] = [];
    const loadedKeys: any[] = [];
    if (data.length > 1) {
      console.log(data, "kkkk");
      const metadata = get(data, "[0]", []);
      const clusters = get(data, "[1]", []);
      const metadataKeys = keys(metadata);
      sortBy(metadataKeys);
      forEach(metadataKeys, (key) => {
        const data = get(metadata, key, {});
        loadedKeys.push(key);
        if (key === StateRoleName.Broker) {
          const clusterNodes: TreeNode[] = [];
          forEach(clusters || [], (cluster: any) => {
            const namespace = cluster.config.namespace;
            const clusterKey = `${key}-${namespace}`;
            const clusterTypes = getItems(clusterKey, key, data, namespace);
            clusterNodes.push({
              label: (
                <>
                  {namespace} (<ClusterStatusView text={cluster.status} />)
                </>
              ),
              value: namespace,
              key: clusterKey,
              parent: key,
              children: clusterTypes,
            });
          });
          root.push({
            label: key,
            value: key,
            key: key,
            children: clusterNodes,
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
    } else {
      const metadata = get(data, "[0]", []);
      const metadataKeys = keys(metadata);
      sortBy(metadataKeys);
      forEach(metadataKeys, (key) => {
        const data = get(metadata, key, {});
        loadedKeys.push(key);
        console.log(key, getItems(key, key, data));
        root.push({
          label: key,
          value: key,
          key: key,
          children: getItems(key, key, data),
        });
      });
    }
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
