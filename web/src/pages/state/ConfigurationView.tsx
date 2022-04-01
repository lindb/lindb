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
import React, { useState } from "react";
import { Card, Descriptions, Space, Typography } from "@douyinfe/semi-ui";
import * as _ from "lodash-es";
import { useWatchURLChange } from "@src/hooks";
import { proxy } from "@src/services";
import { URLStore } from "@src/stores";
import * as monaco from "monaco-editor";

const { Text } = Typography;
/**
 * ConfigurationView which view configuration in node's memory.
 */
export default function ConfigurationView() {
  const [config, setConfig] = useState();
  const [loading, setLoading] = useState(false);
  useWatchURLChange(async () => {
    const target = URLStore.params.get("target");
    if (!target) {
      return;
    }
    setLoading(true);
    try {
      const config = await proxy({
        target: URLStore.params.get("target"),
        path: "/api/config",
      });
      setConfig(config);
    } finally {
      setLoading(false);
    }
  });
  return (
    <>
      <Card bodyStyle={{ padding: 12 }} loading={loading}>
        <Space align="center">
          <Descriptions
            row
            size="small"
            data={[
              {
                key: "Host IP",
                value: (
                  <Text link>{_.get(config, "node.hostIp", "unknown")}</Text>
                ),
              },
              {
                key: "Host Name",
                value: (
                  <Text link>{_.get(config, "node.hostName", "unknown")}</Text>
                ),
              },
              {
                key: "HTTP",
                value: (
                  <Text link>{_.get(config, "node.httpPort", "unknown")}</Text>
                ),
              },
              {
                key: "GRPC",
                value: (
                  <Text link>{_.get(config, "node.grpcPort", "unknown")}</Text>
                ),
              },
            ]}
          />
        </Space>
      </Card>
      <Card
        bodyStyle={{ padding: 0 }}
        style={{ marginTop: 12 }}
        loading={loading}
      >
        {/* <Editor
          language="ini"
          theme="vs-dark"
          value={_.get(config, "config", "")}
          options={{
            readOnly: true,
            minimap: {
              enabled: false,
            },
            padding: { top: 12 },
          }}
          onMount={(editor: any, monaco: any) => {
            const editorElement = editor.getDomNode();

            if (!editorElement) {
              return;
            }

            const lineCount = editor.getModel()?.getLineCount() || 1;
            const height = editor.getTopForLineNumber(lineCount + 1);

            editorElement.style.height = `${height}px`;
            editor.layout();
          }}
        /> */}
      </Card>
    </>
  );
}
