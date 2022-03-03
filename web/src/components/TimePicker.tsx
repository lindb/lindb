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
import {
  Popover,
  Button,
  Col,
  Row,
  Typography,
  Input,
} from "@douyinfe/semi-ui";
import { IconClock } from "@douyinfe/semi-icons";
import { URLStore } from "@src/stores";
import { reaction } from "mobx";
import * as _ from "lodash-es";
import { useWatchURLChange } from "@src/hooks";

const { Text, Title } = Typography;
const defaultQuickItem = { title: "Last 1 hour", value: "now()-1h" };

type QuickSelectItem = {
  title: string;
  value: string;
};

const quickSelectList: QuickSelectItem[] = [
  { title: "Last 15 min", value: "now()-15m" },
  { title: "Last 30 min", value: "now()-30m" },
  defaultQuickItem,
  { title: "Last 3 hours", value: "now()-3h" },
  { title: "Last 6 hours", value: "now()-6h" },
  { title: "Last 12 hours", value: "now()-12h" },
  { title: "Last 1 day", value: "now()-1d" },
  { title: "Last 2 days", value: "now()-2d" },
  { title: "Last 3 days", value: "now()-3d" },
  { title: "Last 7 days", value: "now()-7d" },
  { title: "Last 15 days", value: "now()-15d" },
  { title: "Last 30 days", value: "now()-30d" },
];

export default function TimePicker() {
  const [quick, setQuick] = useState<QuickSelectItem>(defaultQuickItem);
  const [quickItems, setQuickItems] =
    useState<QuickSelectItem[]>(quickSelectList);
  const [visible, setVisible] = useState(false);
  useWatchURLChange(() => {
    const val = URLStore.params.get("from");
    const quickItem = _.find(quickSelectList, { value: `${val}` });
    setQuick(quickItem || defaultQuickItem);
  });

  const renderQuickSelectItem = (
    items: QuickSelectItem[],
    span: number = 12
  ) => {
    const SelectItems = items.map((item) => (
      <div key={item.value} style={{ padding: 4 }}>
        <Text
          link
          onClick={() => {
            setVisible(false);
            URLStore.changeURLParams({
              params: { from: `${item.value}` },
              needDelete: ["from", "to"],
            });
          }}
        >
          {item.title}
        </Text>
      </div>
    ));
    return <Col span={span}>{SelectItems}</Col>;
  };

  /**
   * Render current selected time
   */
  function renderSelectedTime() {
    return (
      <Button icon={<IconClock />} onClick={() => setVisible(true)}>
        {quick.title}
      </Button>
    );
  }

  function renderTimeSelectPanel() {
    return (
      <div style={{ width: 200 }}>
        <Title strong heading={6}>
          <Input
            placeholder="Search quick ranges"
            onChange={(val: string) => {
              const rs = _.filter(
                quickSelectList,
                (item: QuickSelectItem) => item.title.indexOf(val) >= 0
              );
              setQuickItems(rs);
            }}
          />
        </Title>
        <Row>{renderQuickSelectItem(quickItems)}</Row>
      </div>
    );
  }

  return (
    <Popover
      showArrow
      visible={visible}
      trigger="click"
      position="bottom"
      content={renderTimeSelectPanel()}
    >
      {renderSelectedTime()}
    </Popover>
  );
}
