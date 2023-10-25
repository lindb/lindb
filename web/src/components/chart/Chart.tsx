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
import { IconLineChartStroked } from "@douyinfe/semi-icons";
import { Card, Space, Tooltip, Typography } from "@douyinfe/semi-ui";
import { useMetric } from "@src/hooks";
import { ChartConfig, ChartType, Query, ResultSet } from "@src/models";
import { SimpleStatusTip } from "@src/components";
import { ChartKit } from "@src/utils";
import React from "react";
import { LazyLoad } from "..";
import CanvasChart from "./CanvasChart";

const { Text } = Typography;

type ChartProps = {
  type: ChartType;
  queries: Query[];
  config?: ChartConfig;
  height?: number;
  disableBind?: boolean;
};

const InternalChart: React.FC<ChartProps> = (props) => {
  const { type, queries, config, disableBind } = props;
  const { isLoading, isError, data, error } = useMetric(queries || [], {
    disableBind,
  });
  const datasets = ChartKit.createDatasets(data as ResultSet[], type);
  return (
    <Card
      style={{ height: "100%" }}
      bodyStyle={{ padding: 8, height: "calc(100% - 48px)" }}
      headerStyle={{ padding: 6, lineHeight: "16px" }}
      title={
        <Space className="lin-small-space" align="center">
          {config?.description ? (
            <Tooltip content={config.description}>
              <IconLineChartStroked />
            </Tooltip>
          ) : (
            <IconLineChartStroked />
          )}
          <Text>{config?.title}</Text>
        </Space>
      }
      headerExtraContent={
        <Space className="lin-small-space" align="center">
          <SimpleStatusTip
            isLoading={isLoading}
            error={error}
            isError={isError}
          />
        </Space>
      }
    >
      <CanvasChart
        datasets={datasets}
        config={config?.config}
        type={type}
        unit={config?.unit}
      />
    </Card>
  );
};

export const Chart: React.FC<ChartProps> = (props) => {
  const { height } = props;
  return (
    <div style={{ height: height || 300 }}>
      <LazyLoad>
        <InternalChart {...props} />
      </LazyLoad>
    </div>
  );
};
