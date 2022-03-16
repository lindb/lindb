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

import { Badge, Space, Typography } from "@douyinfe/semi-ui";
import React from "react";

const { Text } = Typography;

export default function StorageStatusView(props: {
  text: string;
  showBadge?: boolean;
}) {
  const { text, showBadge } = props;
  let color = "warning";
  switch (text) {
    case "Ready":
      color = "success";
      break;
    case "Initialize":
      color = "secondary";
      break;
  }
  return (
    <Space align="center">
      {showBadge && (
        <Badge
          dot
          style={{
            width: 12,
            height: 12,
            marginTop: 4,
            backgroundColor: `var(--semi-color-${color})`,
          }}
        />
      )}
      <Text style={{ color: `var(--semi-color-${color})` }}> {text}</Text>
    </Space>
  );
}
