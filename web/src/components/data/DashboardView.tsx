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
import { Col, Row } from "@douyinfe/semi-ui";
import { Metric, VariatesSelect } from "@src/components";
import { Dashboard } from "@src/models";
import React from "react";
import { v4 as uuidv4 } from "uuid";

export type DashboardViewProps = {
  dashboard: Dashboard;
};

export default function DashboardView(props: DashboardViewProps) {
  const { dashboard } = props;
  return (
    <>
      {dashboard.variates && <VariatesSelect variates={dashboard.variates} />}
      {dashboard.rows.map((row, rowIdx) => (
        <Row
          style={{ marginTop: 12 }}
          key={rowIdx}
          gutter={dashboard.gutter || 8}
        >
          {row?.panels.map((panel, panelIdx) => (
            <Col span={panel.span || 12} key={`${rowIdx}-${panelIdx}`}>
              <Metric chartId={uuidv4()} config={panel.chart} />
            </Col>
          ))}
        </Row>
      ))}
    </>
  );
}
