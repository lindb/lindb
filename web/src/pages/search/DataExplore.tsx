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
import React, { useEffect } from "react";
import { Card, Form, Popover, Typography, Space } from "@douyinfe/semi-ui";
import { MetadataSelect } from "@src/components";
import { IconFilter, IconHistogram } from "@douyinfe/semi-icons";
import { ChartStore, URLStore } from "@src/stores";
import { CanvasChart } from "@src/components";
const { Text } = Typography;
export default function DataExplore() {
  useEffect(() => {
    const chartId = "xxxxxxxx";
    // register chart config
    ChartStore.register(chartId, {
      title: "CPU",
      targets: [
        {
          db: "_internal",
          ql: "explain select used_percent from lindb.monitor.system.disk_usage_stats",
        },
      ],
    });
    URLStore.forceChange(); // trigger url change
    return () => {
      // unRegister chart config when component destroy.
      ChartStore.unRegister(chartId);
      // console.log("unregister chart");
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  return (
    <>
      <Card
        bordered={false}
        style={{ marginBottom: 12 }}
        bodyStyle={{ padding: 12 }}
      >
        <Form layout="horizontal">
          <MetadataSelect
            variate={{ tagKey: "db", label: "Database", ql: "show databases" }}
            labelPosition="inset"
          />
          <MetadataSelect
            variate={{
              db: "_internal",
              tagKey: "namespace",
              label: "Namespace",
              ql: "show namespaces",
            }}
            labelPosition="inset"
          />
          <MetadataSelect
            variate={{
              db: "_internal",
              tagKey: "namespace",
              label: "Metrics",
              ql: "show metrics",
            }}
            labelPosition="inset"
          />
        </Form>
      </Card>
      <Card
        title={
          <Space align="center">
            <IconHistogram />
            <Text strong>cpu.usage</Text>
          </Space>
        }
        headerStyle={{ padding: 12 }}
        bordered={false}
        style={{ marginBottom: 12 }}
        bodyStyle={{ padding: 12 }}
      >
        <Form
          layout="horizontal"
          getFormApi={(api) => {
            setTimeout(() => {
              api.setValues({
                groupBy: ["node", "role"],
                field: ["sys", "steal"],
              });
            }, 1000);
          }}
        >
          <Form.Select
            field="groupBy"
            labelPosition="left"
            label="Group By:"
            multiple
            optionList={[
              { value: "node", label: "Node" },
              { value: "role", label: "role" },
            ]}
          />
          <Form.Slot labelPosition="left" label="Filter By:">
            <div
              style={{ display: "flex", alignItems: "center", height: "100%" }}
            >
              <Popover trigger="click" content={<></>}>
                <IconFilter />
              </Popover>
            </div>
          </Form.Slot>
          <Form.Select
            field="field"
            labelPosition="left"
            label="Field:"
            multiple
            optionList={[
              { value: "sys", label: "sys" },
              { value: "steal", label: "steal" },
            ]}
          />
        </Form>
        <CanvasChart chartId="xxxxxxxx" />
      </Card>
    </>
  );
}
