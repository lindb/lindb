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

import { IconHistogram } from "@douyinfe/semi-icons";
import { Card, Space, Typography, Tooltip } from "@douyinfe/semi-ui";
import { CanvasChart, MetricStatus } from "@src/components";
import { ChartConfig } from "@src/models";
import { ChartStore, URLStore } from "@src/stores";
import React, { useEffect } from "react";
import * as _ from "lodash-es";
const { Text } = Typography;

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
  }, [chartId]);
  return (
    <Card
      title={
        <Space align="center">
          <Space>
            {config.description ? (
              <Tooltip content={config.description}>
                <IconHistogram />
              </Tooltip>
            ) : (
              <IconHistogram />
            )}

            <Text strong>{config.title}</Text>
          </Space>
        </Space>
      }
      headerExtraContent={<MetricStatus chartId={chartId} />}
      // bodyStyle={{ height: "100vh" }}
      headerStyle={{ padding: 12 }}
    >
      <CanvasChart chartId={chartId} />
    </Card>
  );
}
