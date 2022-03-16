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
import {
  IconHelpCircleStroked,
  IconPlay,
  IconSend,
} from "@douyinfe/semi-icons";
import { Button, Card, Col, Form, Row, Space } from "@douyinfe/semi-ui";
import { CanvasChart, MetadataSelect } from "@src/components";
import { SQL } from "@src/constants";
import { useWatchURLChange } from "@src/hooks";
import { ChartStore, URLStore } from "@src/stores";
import React, { MutableRefObject, useEffect, useRef } from "react";

const chartID = "9999999999999999";

export default function DataSearch() {
  const formApi = useRef() as MutableRefObject<any>;
  const query = () => {
    const target = formApi.current.getValues();
    ChartStore.reRegister(chartID, { targets: [target] });
    URLStore.changeURLParams({ params: target, forceChange: true });
  };
  useWatchURLChange(() => {
    if (formApi.current) {
      formApi.current.setValues({
        db: URLStore.params.get("db"),
        sql: URLStore.params.get("sql"),
      });
    }
  });
  useEffect(() => {
    ChartStore.register(chartID, {
      targets: [
        {
          db: URLStore.params.get("db") || "",
          sql: URLStore.params.get("sql") || "",
        },
      ],
    });
    return () => {
      // unRegister chart config when component destroy.
      ChartStore.unRegister(chartID);
    };
  }, []);
  return (
    <>
      <Card
        bordered={false}
        style={{ marginBottom: 12 }}
        bodyStyle={{ padding: 12 }}
      >
        <Form
          style={{ paddingTop: 0, paddingBottom: 0 }}
          getFormApi={(api) => (formApi.current = api)}
        >
          <MetadataSelect
            type="db"
            variate={{
              tagKey: "db",
              label: "Database",
              sql: SQL.ShowDatabases,
            }}
            labelPosition="left"
          />
          <Row>
            <Col span={24}>
              <Form.TextArea
                showClear
                field="sql"
                initValue={URLStore.params.get("sql")}
                label={
                  <Space align="center">
                    <span>LinDB Query Language</span>
                    <IconHelpCircleStroked />
                  </Space>
                }
              />
            </Col>
          </Row>
          <Button
            style={{ marginRight: 12 }}
            icon={<IconPlay size="large" />}
            onClick={query}
          >
            Search
          </Button>
          <Button type="secondary" icon={<IconSend size="large" />}>
            Explain
          </Button>
        </Form>
      </Card>
      <Card bordered={false} style={{ marginTop: 12 }}>
        <CanvasChart chartId={chartID} height={300} />
      </Card>
    </>
  );
}
