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
import { Card, Empty, Form, Typography } from "@douyinfe/semi-ui";
import { Icon, LinSelect, StatusTip } from "@src/components";
import { StateRoleName, SQL } from "@src/constants";
import { UIContext } from "@src/context/UIContextProvider";
import { useParams } from "@src/hooks";
import { Unit } from "@src/models";
import { ExecService, ProxyService } from "@src/services";
import { URLStore } from "@src/stores";
import { FormatKit } from "@src/utils";
import { useQuery } from "@tanstack/react-query";
import * as _ from "lodash-es";
import React, { useContext } from "react";

const { Text } = Typography;

const getLogColor = (text: string) => {
  switch (text) {
    case "INFO":
    case "POST":
      return "success";
    case "ERROR":
    case "DELETE":
      return "danger";
    case "WARN":
    case "PUT":
      return "warning";
    case "DEBUG":
    case "GET":
      return "link";
  }
};

const LogContent: React.FC = () => {
  const { node, file, size } = useParams(["node", "file", "size"]);
  const { locale } = useContext(UIContext);
  const { Common } = locale;
  const { isError, error, data, isInitialLoading } = useQuery(
    ["tail_log", node, file, size],
    async () => {
      const renderLogs = (text: string) => {
        if (!text) {
          return null;
        }
        const textArray = text.split(
          /(INFO|DEBUG|ERROR|WARN|POST|PUT|GET|DELETE)/g
        );
        return textArray.map((str, idx) => {
          return (
            <Text
              key={idx}
              style={{ color: `var(--semi-color-${getLogColor(str)})` }}
              strong
            >
              {str}
            </Text>
          );
        });
      };
      return ProxyService.proxy({
        target: node,
        file: file,
        size: size || 256 * 1024,
        path: "/api/v1/log/view",
      }).then((data) => renderLogs(data));
    },
    {
      enabled: !_.isEmpty(node) && !_.isEmpty(file),
    }
  );

  if (isInitialLoading || isError) {
    return (
      <StatusTip
        style={{ marginTop: 100, marginBottom: 100 }}
        isLoading={isInitialLoading}
        isError={isError}
        error={error}
      />
    );
  }
  if (!data) {
    return (
      <Empty
        image={<Icon icon="iconempty" style={{ fontSize: 48 }} />}
        description={Common.noData}
        style={{ marginTop: 100, marginBottom: 100 }}
      />
    );
  }
  return (
    <pre style={{ wordWrap: "normal", whiteSpace: "pre", overflow: "auto" }}>
      {data}
    </pre>
  );
};

const LogView: React.FC = () => {
  const { locale, env } = useContext(UIContext);
  const { LogView } = locale;
  return (
    <>
      <Card style={{ marginBottom: 12 }} bodyStyle={{ padding: 12 }}>
        <Form layout="horizontal" style={{ paddingTop: 0, paddingBottom: 0 }}>
          <LinSelect
            field="role"
            label={LogView.role}
            loader={() =>
              env.role === StateRoleName.Broker
                ? [
                    {
                      value: StateRoleName.Broker,
                      label: StateRoleName.Broker,
                    },
                    {
                      value: StateRoleName.Storage,
                      label: StateRoleName.Storage,
                    },
                  ]
                : [{ value: StateRoleName.Root, label: StateRoleName.Root }]
            }
            style={{ width: 150 }}
            clearKeys={["storage", "node", "file"]}
          />

          <LinSelect
            field="storage"
            label={LogView.storage}
            style={{ width: 200 }}
            loader={() =>
              ExecService.exec<any[]>({ sql: SQL.ShowStorageAliveNodes }).then(
                (data) =>
                  _.map(data || [], (s) => {
                    return { value: s.name, label: s.name };
                  })
              )
            }
            visible={() =>
              _.get(URLStore.getParams(), "role") === StateRoleName.Storage
            }
            reloadKeys={["role"]}
            clearKeys={["node", "file"]}
          />
          <LinSelect
            field="node"
            label={LogView.node}
            style={{ width: 230 }}
            loader={async () => {
              const params = URLStore.getParams();
              const role = _.get(params, "role");
              if (
                role === StateRoleName.Broker ||
                role === StateRoleName.Root
              ) {
                return ExecService.exec<any>({
                  sql:
                    role === StateRoleName.Broker
                      ? SQL.ShowBrokerAliveNodes
                      : SQL.ShowRootAliveNodes,
                }).then((data) =>
                  _.map(data || [], (n: any) => {
                    const target = `${n.hostIp}:${n.httpPort}`;
                    return { value: target, label: target };
                  })
                );
              } else {
                return ExecService.exec<any[]>({
                  sql: SQL.ShowStorageAliveNodes,
                }).then((data) => {
                  const nodes = _.get(
                    _.find(data, {
                      name: _.get(params, "storage"),
                    }),
                    "liveNodes",
                    []
                  );
                  return _.map(nodes, (n: any) => {
                    const target = `${n.hostIp}:${n.httpPort}`;
                    return { value: target, label: target };
                  });
                });
              }
            }}
            clearKeys={["file"]}
            reloadKeys={["role", "storage"]}
          />
          <LinSelect
            field="file"
            label={LogView.file}
            loader={() => {
              const params = URLStore.getParams();
              const target = _.get(params, "node");
              if (!target) {
                return null;
              }
              return ProxyService.proxy({
                target: target,
                path: "/api/v1/log/list",
              }).then((files) =>
                _.map(files || [], (f: any) => {
                  return {
                    value: f.name,
                    label: `${f.name}(${FormatKit.format(f.size, Unit.Bytes)})`,
                  };
                })
              );
            }}
            style={{ width: 240 }}
            reloadKeys={["node"]}
          />
          <LinSelect
            field="size"
            label={LogView.size}
            loader={() => [
              { label: "256KB", value: `${256 * 1024}` },
              { label: "1MB", value: `${1024 * 1024}` },
              { label: "3MB", value: `${3 * 1024 * 1024}` },
              { label: "5MB", value: `${5 * 1024 * 1024}` },
            ]}
            style={{ width: 140 }}
          />
        </Form>
      </Card>
      <Card bodyStyle={{ padding: 12 }}>
        <LogContent />
      </Card>
    </>
  );
};

export default LogView;
