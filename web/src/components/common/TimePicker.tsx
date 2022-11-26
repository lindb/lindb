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
import { IconChevronDown, IconRefresh, IconTick } from "@douyinfe/semi-icons";
import {
  Button,
  Dropdown,
  Input,
  Popover,
  Form,
  SplitButtonGroup,
  Typography,
  Space,
} from "@douyinfe/semi-ui";
import { useParams } from "@src/hooks";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, {
  useState,
  useRef,
  MutableRefObject,
  useEffect,
  useContext,
} from "react";
import moment from "moment";
import { DateTimeFormat } from "@src/constants";
import { Icon } from "@src/components";
import { UIContext } from "@src/context/UIContextProvider";

const { Title } = Typography;
type QuickSelectItem = {
  title: string;
  value: string;
};
const defaultQuickItem = { title: "last1Hour", value: "now()-1h" };
const defaultAutoRefreshItem = { title: "off", value: "" };
const quickSelectList: QuickSelectItem[] = [
  { title: "last15Min", value: "now()-15m" },
  { title: "last30Min", value: "now()-30m" },
  defaultQuickItem,
  { title: "last3Hour", value: "now()-3h" },
  { title: "last6Hour", value: "now()-6h" },
  { title: "last12Hour", value: "now()-12h" },
  { title: "last1Day", value: "now()-1d" },
  { title: "last2Day", value: "now()-2d" },
  { title: "last3Day", value: "now()-3d" },
  { title: "last7Day", value: "now()-7d" },
  { title: "last15Day", value: "now()-15d" },
  { title: "klast30Day", value: "now()-30d" },
];
const autoRefreshList: QuickSelectItem[] = [
  {
    value: "",
    title: "off",
  },
  { value: "10", title: `10s` },
  { value: "30", title: `30s` },
  { value: "60", title: `1m` },
  { value: "300", title: `5m` },
];

export default function TimePicker() {
  const { from, to, refresh } = useParams(["from", "to", "refresh"]);
  const { locale } = useContext(UIContext);
  const { TimePicker } = locale;
  const formApi = useRef() as MutableRefObject<any>;
  const [quick, setQuick] = useState<QuickSelectItem | undefined>(
    defaultQuickItem
  );
  const [quickItems, setQuickItems] =
    useState<QuickSelectItem[]>(quickSelectList);
  const [visible, setVisible] = useState(false);
  const [autoRefresh, setAutoRefresh] = useState<QuickSelectItem>(
    defaultAutoRefreshItem
  );
  const countDown = useRef<number>();
  const timeRangeVisible = useRef<boolean>(false);

  const buildCountDown = (interval: number) => {
    if (countDown.current) {
      clearInterval(countDown.current);
    }

    if (interval) {
      countDown.current = +setInterval(() => {
        URLStore.forceChange();
      }, 1000 * interval);
    }
  };

  useEffect(() => {
    if (_.isEmpty(from)) {
      setQuick(defaultQuickItem);
    } else {
      const quickItem = _.find(quickSelectList, { value: `${from}` });
      setQuick(quickItem);
    }
  }, [from]);

  useEffect(() => {
    const refreshItem = _.find(autoRefreshList, { value: `${refresh}` });
    if (refreshItem && refreshItem.value !== "") {
      buildCountDown(parseInt(refreshItem.value));
    } else {
      clearInterval(countDown.current);
    }
    setAutoRefresh(refreshItem || defaultAutoRefreshItem);
  }, [refresh]);

  const renderQuickSelectItem = (items: QuickSelectItem[]) => {
    const SelectItems = items.map((item) => (
      <Dropdown.Item
        style={{ padding: 3 }}
        key={item.title}
        active={quick?.value == item.value}
        onClick={() => {
          setVisible(false);
          URLStore.changeURLParams({
            params: { from: `${item.value}` },
            needDelete: ["from", "to"],
          });
        }}
      >
        <IconTick
          style={{
            color: quick?.value !== item.value ? "transparent" : "inherit",
          }}
        />
        {TimePicker[item.title]}
      </Dropdown.Item>
    ));
    return <Dropdown.Menu>{SelectItems}</Dropdown.Menu>;
  };

  /**
   * Render current selected time
   */
  function renderSelectedTime() {
    return (
      <Button icon={<Icon icon="iconclock" />} onClick={() => setVisible(true)}>
        {quick && TimePicker[quick.title]}
        {!quick && `${from} ~ ${to ? `${to}` : "now"}`}
      </Button>
    );
  }

  function renderTimeSelectPanel() {
    return (
      <Space style={{ width: 460, padding: 20 }} align="start">
        <div style={{ width: 230 }}>
          <Title heading={5}>{TimePicker.absoluteTimeRange}</Title>
          <Form
            style={{ marginTop: 16 }}
            className="lin-form"
            getFormApi={(api: any) => (formApi.current = api)}
          >
            <Form.DatePicker
              field="from"
              type="dateTime"
              label={TimePicker.from}
              labelPosition="top"
              onOpenChange={(v) => (timeRangeVisible.current = v)}
              initValue={!quick ? from && new Date(`${from}`) : null}
            />
            <Form.DatePicker
              field="to"
              type="dateTime"
              labelPosition="top"
              label={TimePicker.to}
              onOpenChange={(v) => (timeRangeVisible.current = v)}
              initValue={to && new Date(`${to}`)}
            />
            <Button
              style={{ marginTop: 12 }}
              onClick={() => {
                setVisible(false);
                const from = formApi.current.getValue("from");
                const to = formApi.current.getValue("to");
                URLStore.changeURLParams({
                  params: {
                    from: from
                      ? moment(from.getTime()).format(DateTimeFormat)
                      : "",
                    to: to ? moment(to.getTime()).format(DateTimeFormat) : "",
                  },
                });
              }}
            >
              {TimePicker.applyTimeRange}
            </Button>
          </Form>
        </div>
        <div
          style={{
            paddingLeft: 20,
            borderLeft: "1px solid var(--semi-color-text-3)",
          }}
        >
          <Title strong heading={6}>
            <Input
              placeholder={TimePicker.searchQuickRange}
              onChange={(val: string) => {
                const rs = _.filter(
                  quickSelectList,
                  (item: QuickSelectItem) => item.title.indexOf(val) >= 0
                );
                setQuickItems(rs);
              }}
            />
          </Title>
          {renderQuickSelectItem(quickItems)}
        </div>
      </Space>
    );
  }

  return (
    <>
      <Popover
        onClickOutSide={(_v) => {
          if (!timeRangeVisible.current) {
            // if click outside not date time range picker
            setVisible(false);
          }
        }}
        showArrow
        visible={visible}
        trigger="custom"
        position="bottom"
        content={renderTimeSelectPanel()}
      >
        {renderSelectedTime()}
      </Popover>
      <SplitButtonGroup style={{ marginLeft: 8 }}>
        <Button icon={<IconRefresh />} onClick={() => URLStore.forceChange()} />
        <Dropdown
          trigger="click"
          showTick
          render={
            <Dropdown.Menu>
              {autoRefreshList.map((item) => (
                <Dropdown.Item
                  key={item.title}
                  active={item.value === autoRefresh.value}
                  onClick={() => {
                    URLStore.changeURLParams({
                      params: { refresh: item.value },
                    });
                  }}
                >
                  {TimePicker[item.title]}
                </Dropdown.Item>
              ))}
            </Dropdown.Menu>
          }
        >
          {autoRefresh.value === "" ? (
            <Button icon={<IconChevronDown />} iconPosition="right" />
          ) : (
            <Button icon={<IconChevronDown />} iconPosition="right">
              {TimePicker[autoRefresh.title]}
            </Button>
          )}
        </Dropdown>
      </SplitButtonGroup>
    </>
  );
}
