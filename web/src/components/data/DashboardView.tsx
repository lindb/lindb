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
import { IconGridStroked } from "@douyinfe/semi-icons";
import { Card, Col, Row, Form, Select } from "@douyinfe/semi-ui";
import { Metric, VariatesSelect } from "@src/components";
import { Dashboard, DashboardItem } from "@src/models";
import { URLStore } from "@src/stores";
import { useWatchURLChange } from "@src/hooks";
import React, { useState } from "react";
import * as _ from "lodash-es";
import { v4 as uuidv4 } from "uuid";

export type DashboardViewProps = {
  dashboards?: DashboardItem[];
};

export default function DashboardView(props: DashboardViewProps) {
  const { dashboards } = props;
  const [selectedDashboard, setSelectedDashboard] = useState(
    _.get(dashboards, "[0].value", "")
  );
  const [dashboard, setDashboard] = useState<Dashboard>(
    _.get(dashboards, "[0].dashboard", {})
  );

  const changeDashboard = (value: string) => {
    const dashboardItem = _.find(dashboards, (item: DashboardItem) => {
      return item.value == value;
    });
    setDashboard(_.get(dashboardItem, "dashboard", {}) as Dashboard);
    setSelectedDashboard(value);
  };

  useWatchURLChange(() => {
    const d = URLStore.params.get("d");
    if (d) {
      changeDashboard(d as string);
    }
  });

  return (
    <>
      <Card>
        {dashboards && dashboards.length > 1 && (
          <Select
            prefix={<IconGridStroked />}
            value={selectedDashboard}
            optionList={dashboards}
            onChange={(value) => {
              changeDashboard(value as string);
              URLStore.changeURLParams({ params: { d: value } });
            }}
            style={{ minWidth: 60, marginRight: 16 }}
          />
        )}
        {dashboard.variates && <VariatesSelect variates={dashboard.variates} />}
      </Card>
      {(dashboard.rows || []).map((row, rowIdx) => (
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
