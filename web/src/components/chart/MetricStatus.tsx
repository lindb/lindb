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
import { IconSpin, IconUploadError } from "@douyinfe/semi-icons";
import { Tooltip, Typography, Banner } from "@douyinfe/semi-ui";
import { ChartStatus } from "@src/models";
import { ChartStore } from "@src/stores";
import { reaction } from "mobx";
import React, { useEffect, useState } from "react";
const { Text } = Typography;

/**
 *  Render metric chart render status.
 *
 * @param props chartId which need watch
 */
export default function MetricStatus(props: {
  chartId: string;
  showMsg?: boolean;
}) {
  const { chartId, showMsg } = props;
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
      if (showMsg) {
        return (
          <Banner
            type="danger"
            description={ChartStore.chartErrMap.get(chartId)}
            closeIcon
            fullMode={false}
          />
        );
      }
      return (
        <Tooltip position="left" content={ChartStore.chartErrMap.get(chartId)}>
          <IconUploadError style={{ color: "var(--semi-color-danger)" }} />
        </Tooltip>
      );
    default:
      return <></>;
  }
}
