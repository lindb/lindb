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
import { ChartStore, URLStore } from "@src/stores";
import { ChartConfig } from "@src/models";
import { CanvasChart } from "@src/components";
import { Card, Space, Typography, Tooltip } from "@douyinfe/semi-ui";
import { IconHistogram, IconSpin, IconUploadError } from "@douyinfe/semi-icons";
import { reaction } from "mobx";
import { ChartStatus } from "@src/models";
const { Text } = Typography;

/**
 *  Render metric chart render status.
 *
 * @param props chartId which need watch
 */
const MetricStatus = (props: { chartId: string }) => {
  const { chartId } = props;
  const [status, setStatus] = useState(ChartStatus.Init);
  useEffect(() => {
    const disposer = [
      reaction(
        () => ChartStore.chartStatusMap.get(chartId),
        (s: ChartStatus | undefined) => {
          // watch chart status, if change set state.
          if (s) {
            setStatus(s);
          }
        }
      ),
    ];
    return () => {
      disposer.forEach((d) => d());
    };
  }, [chartId]);

  switch (status) {
    case ChartStatus.Loading:
      return <IconSpin spin style={{ color: "var(--semi-color-primary)" }} />;
    case ChartStatus.Error:
      return (
        <Tooltip position="left" content={ChartStore.chartErrMap.get(chartId)}>
          <IconUploadError style={{ color: "var(--semi-color-danger)" }} />
        </Tooltip>
      );
    default:
      return <></>;
  }
};

interface MetricProps {
  chartId: string;
  config: ChartConfig;
}
export default function Metric(props: MetricProps) {
  const { chartId, config } = props;
  useEffect(() => {
    // register chart config
    ChartStore.register(chartId, config);
    URLStore.forceChange(); // trigger url change
    return () => {
      // unRegister chart config when component destroy.
      ChartStore.unRegister(chartId);
      // console.log("unregister chart");
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  return (
    <Card
      title={
        <Space align="center">
          <IconHistogram />
          <Text strong>{config.title}</Text>
        </Space>
      }
      headerExtraContent={<MetricStatus chartId={chartId} />}
      bordered={false}
      // bodyStyle={{ height: "100vh" }}
      headerStyle={{ padding: 12 }}
    >
      <CanvasChart chartId={chartId} />
    </Card>
  );
}
