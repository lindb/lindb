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
import React from "react";
import { Progress, Descriptions, Typography } from "@douyinfe/semi-ui";
import { FormatKit } from "@src/utils";
import { Unit } from "@src/models";

const { Text } = Typography;

interface CapacityViewProps {
  percent: number;
  total: number;
  free: number;
  used: number;
}

export default function CapacityView(props: CapacityViewProps) {
  const { percent, total, free, used } = props;
  return (
    <>
      <div style={{ width: "70%", marginBottom: 8 }}>
        <Progress
          motion={false}
          className="lin-stats"
          percent={percent}
          stroke="#fc8800"
          size="large"
          format={(val) => FormatKit.format(val, Unit.Percent)}
          showInfo={true}
        />
      </div>
      <Descriptions
        className="lin-small-desc"
        row
        size="small"
        data={[
          {
            key: "Total",
            value: <Text link>{FormatKit.format(total, Unit.Bytes)}</Text>,
          },
          {
            key: "Used",
            value: (
              <Text type="warning">{FormatKit.format(used, Unit.Bytes)}</Text>
            ),
          },
          {
            key: "Free",
            value: (
              <Text type="success">{FormatKit.format(free, Unit.Bytes)}</Text>
            ),
          },
        ]}
      />
    </>
  );
}
