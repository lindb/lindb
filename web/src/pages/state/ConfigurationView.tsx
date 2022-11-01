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
import React, { useEffect, useRef, MutableRefObject, useContext } from "react";
import { Card, Descriptions, Space, Typography } from "@douyinfe/semi-ui";
import * as _ from "lodash-es";
import { useParams } from "@src/hooks";
import { ProxyService } from "@src/services";
import * as monaco from "monaco-editor";
import { useQuery } from "@tanstack/react-query";
import { StatusTip } from "@src/components";
import { UIContext } from "@src/context/UIContextProvider";
import { Theme } from "@src/constants";

const { Text } = Typography;
/**
 * ConfigurationView which view configuration in node's memory.
 */
const ConfigurationView: React.FC = () => {
  const editorRef = useRef() as MutableRefObject<HTMLDivElement>;
  const { target } = useParams(["target"]);
  const { theme } = useContext(UIContext);
  const {
    isLoading,
    isInitialLoading,
    isFetching,
    isError,
    data: config,
    error,
  } = useQuery(
    ["show_cfg", target],
    async () => {
      return ProxyService.proxy({
        target: target,
        path: "/api/v1/config",
      });
    },
    {
      enabled: !_.isEmpty(target),
    }
  );

  useEffect(() => {
    if (isLoading || isError || !editorRef.current) {
      return;
    }
    monaco.editor.create(editorRef.current, {
      language: "ini",
      theme: theme === Theme.dark ? "vs-dark" : "vs",
      // lineNumbers: "off",
      minimap: { enabled: false },
      // lineNumbersMinChars: 2,
      readOnly: true,
      value: _.get(config, "config", "no data"),
    });
  }, [config, isLoading, isError, editorRef, theme]);

  if (isError || isLoading || isInitialLoading || isFetching) {
    return (
      <StatusTip
        style={{ marginTop: 150 }}
        isLoading={isLoading || isInitialLoading || isFetching}
        isError={isError}
        error={error}
      />
    );
  }

  return (
    <>
      <Card bodyStyle={{ padding: 12 }}>
        <Space align="center">
          <Descriptions
            row
            size="small"
            className="lin-description"
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
      <Card bodyStyle={{ padding: 0 }} style={{ marginTop: 12 }}>
        <div ref={editorRef} style={{ height: "90vh" }} />
      </Card>
    </>
  );
};

export default ConfigurationView;
