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
import React, { useState, useEffect } from "react";
import { MetadataSelect, ExplainResultView } from "@src/components";
import { queryMetric } from "@src/services";
import { ResultSet } from "@src/models";
import { ChartStore } from "@src/stores";
import Metric from "@src/components/chart/Metric";

const chartID = "9999999999999999";

export default function DataSearch() {
  const [resultSet, setResultSet] = useState<ResultSet>();
  const query = async () => {
    const rs = await queryMetric({
      db: "_internal",
      sql: "explain select used_percent from lindb.monitor.system.disk_usage_stats",
    });
    setResultSet(rs);
  };
  useEffect(() => {
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
        <Form style={{ paddingTop: 0, paddingBottom: 0 }}>
          <MetadataSelect
            variate={{ tagKey: "db", label: "Database", ql: "show databases" }}
            labelPosition="left"
          />
          <Row>
            <Col span={24}>
              <Form.TextArea
                showClear
                field="ql"
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
      <Metric
        chartId={chartID}
        config={{
          title: "CPU",
          targets: [
            {
              db: "_internal",
              ql: "explain select used_percent from lindb.monitor.system.disk_usage_stats",
            },
          ],
        }}
      />
      <Card bordered={false} style={{ marginTop: 12 }}>
        {resultSet?.stats && <ExplainResultView stats={resultSet?.stats} />}
      </Card>
    </>
  );
}
