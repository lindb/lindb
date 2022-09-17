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
  IllustrationIdle,
  IllustrationIdleDark,
} from "@douyinfe/semi-illustrations";
import { Card, Empty, Form, Typography, useFormState } from "@douyinfe/semi-ui";
import { StateRoleName, SQL } from "@src/constants";
import { useAliveState, useStorage } from "@src/hooks";
import { UnitEnum } from "@src/models";
import { proxy } from "@src/services";
import { URLStore } from "@src/stores";
import { formatter } from "@src/utils";
import * as _ from "lodash-es";
import React, { MutableRefObject, useRef, useState } from "react";

const { Text } = Typography;

export default function LogView() {
  const { aliveState: liveNodes } = useAliveState(SQL.ShowBrokerAliveNodes);
  const { storages } = useStorage();
  const [tailing, setTailing] = useState(false);
  const [loading, setLoading] = useState(false);
  const [logs, setLogs] = useState<any>(null);
  const [files, setFiles] = useState([]);
  const [nodes, setNodes] = useState<any[]>([]);
  const formApi = useRef() as MutableRefObject<any>;

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

  const tailLog = async (params: any) => {
    if (!_.get(params, "target") || !_.get(params, "file")) {
      return;
    }
    URLStore.changeURLParams({
      params: formApi.current.getValues(),
      clearAll: true,
    });

    setTailing(true);
    try {
      const logs = await proxy({ ...params, path: "/api/v1/log/view" });

      setLogs(renderLogs(logs as string));
    } finally {
      setTailing(false);
    }
  };

  const listFiles = async (target: any) => {
    setLoading(true);
    try {
      const files = await proxy({
        target: target,
        path: "/api/v1/log/list",
      });
      setFiles(files);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectRole = (_value: any, _option: any) => {
    formApi.current.setValues({ node: "", file: "", storage: "" });
  };

  const StorageSelectInput = () => {
    const formState = useFormState();

    return (
      <>
        {_.get(formState, "values.role") === StateRoleName.Storage && (
          <Form.Select
            field="storage"
            placeholder={`Please select`}
            optionList={_.map(storages || [], (s) => {
              return { value: s.name, label: s.name };
            })}
            labelPosition="inset"
            label="Storage"
            style={{ width: 180 }}
          />
        )}
      </>
    );
  };

  return (
    <>
      <Card style={{ marginBottom: 12 }} bodyStyle={{ padding: 12 }}>
        <Form
          getFormApi={(api) => {
            formApi.current = api;
            URLStore.params.forEach((value: string, key: string) => {
              formApi.current.setValue(key, value);
            });
          }}
          onValueChange={(values: any) => {
            const params = _.pick(values, ["target", "file", "size"]);
            tailLog({ ...params });
          }}
          layout="horizontal"
          style={{ paddingTop: 0, paddingBottom: 0 }}
        >
          <Form.Select
            field="role"
            placeholder={`Please select`}
            optionList={[
              { value: StateRoleName.Broker, label: StateRoleName.Broker },
              {
                value: StateRoleName.Storage,
                label: StateRoleName.Storage,
              },
            ]}
            labelPosition="inset"
            label="Role"
            style={{ width: 150 }}
            onSelect={handleSelectRole}
          />
          <StorageSelectInput />
          <Form.Select
            field="target"
            rules={[{ required: true }]}
            placeholder={`Please select`}
            optionList={nodes}
            labelPosition="inset"
            label="Node"
            loading={loading}
            style={{ width: 230 }}
            onFocus={() => {
              if (formApi.current.getValue("role") === StateRoleName.Broker) {
                setNodes(
                  _.map(liveNodes || [], (n: any) => {
                    const target = `${n.hostIp}:${n.httpPort}`;
                    return { value: target, label: target };
                  })
                );
              } else {
                const nodes = _.get(
                  _.find(storages, {
                    name: formApi.current.getValue("storage"),
                  }),
                  "liveNodes",
                  []
                );
                setNodes(
                  _.map(nodes, (n: any) => {
                    const target = `${n.hostIp}:${n.httpPort}`;
                    return { value: target, label: target };
                  })
                );
              }
            }}
          />
          <Form.Select
            field="file"
            placeholder={`Please select`}
            optionList={_.map(files || [], (f: any) => {
              return {
                value: f.name,
                label: `${f.name}(${formatter(f.size, UnitEnum.Bytes)})`,
              };
            })}
            labelPosition="inset"
            label="File"
            style={{ width: 240 }}
            loading={loading}
            onFocus={() => {
              listFiles(formApi.current.getValue("target"));
            }}
          />
          <Form.Select
            field="size"
            placeholder={`Please select`}
            optionList={[
              { label: "256KB", value: 256 * 1024 },
              { label: "1MB", value: 1024 * 1024 },
              { label: "3MB", value: 3 * 1024 * 1024 },
              { label: "5MB", value: 5 * 1024 * 1024 },
            ]}
            labelPosition="inset"
            label="Last"
            style={{ width: 140 }}
          />
        </Form>
      </Card>
      <Card bodyStyle={{ padding: 12 }} loading={tailing}>
        {!logs ? (
          <Empty
            image={<IllustrationIdle style={{ width: 150, height: 150 }} />}
            darkModeImage={
              <IllustrationIdleDark style={{ width: 150, height: 150 }} />
            }
            description="No log data"
            style={{ marginTop: 50, minHeight: 400 }}
          />
        ) : (
          <pre
            style={{ wordWrap: "normal", whiteSpace: "pre", overflow: "auto" }}
          >
            {logs}
          </pre>
        )}
      </Card>
    </>
  );
}
