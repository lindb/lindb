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
import { IconFixedStroked, IconGridStroked } from "@douyinfe/semi-icons";
import { Card, Col, Row, Select } from "@douyinfe/semi-ui";
import { Metric, VariatesSelect } from "@src/components";
import { StateRoleName } from "@src/constants";
import { useWatchURLChange } from "@src/hooks";
import { Dashboard, DashboardItem } from "@src/models";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { useState, useEffect } from "react";
import { v4 as uuidv4 } from "uuid";

export type DashboardViewProps = {
  dashboards?: DashboardItem[];
  variates?: any;
};

export default function DashboardView(props: DashboardViewProps) {
  const { dashboards, variates } = props;
  const [selectedDashboard, setSelectedDashboard] = useState(
    _.get(dashboards, "[0]", {})
  );
  const [dashboard, setDashboard] = useState<Dashboard>(
    _.get(dashboards, "[0].dashboard", {})
  );
  const [role, setRole] = useState<any>(null);

  useEffect(() => {
    URLStore.changeURLParams({
      params: { role: URLStore.params.get("role") || StateRoleName.Broker },
    });
  }, []);

  const changeDashboard = (value: string) => {
    const r = URLStore.params.get("role") || StateRoleName.Broker;
    const currentDashboards = _.filter(
      dashboards,
      (o) => _.indexOf(o.scope, r) >= 0
    );
    let dashboardItem = _.find(currentDashboards, (item: DashboardItem) => {
      return item.value == value;
    });
    if (_.isUndefined(dashboardItem) || _.isNull(dashboardItem)) {
      dashboardItem = currentDashboards[0];
    }
    setDashboard(_.get(dashboardItem, "dashboard", {}) as Dashboard);
    setSelectedDashboard(dashboardItem);
  };

  useWatchURLChange(() => {
    const role = URLStore.params.get("role") || StateRoleName.Broker;
    setRole(role);
    const d = URLStore.params.get("d");
    changeDashboard(d as string);
  });

  return (
    <>
      <Card bodyStyle={{ padding: 12 }}>
        <Select
          value={role}
          prefix={<IconFixedStroked />}
          optionList={[
            { value: StateRoleName.Broker, label: StateRoleName.Broker },
            { value: StateRoleName.Storage, label: StateRoleName.Storage },
          ]}
          onChange={(value) => {
            URLStore.changeURLParams({
              params: { role: value },
              needDelete: ["namespace", "node"],
            });
          }}
          style={{ minWidth: 60, marginRight: 16 }}
        />
        <Select
          value={selectedDashboard?.value}
          prefix={<IconGridStroked />}
          onChange={(value) => {
            URLStore.changeURLParams({ params: { d: value } });
          }}
          optionList={_.filter(
            dashboards,
            (o) => _.indexOf(o.scope, role) >= 0
          )}
          style={{ minWidth: 60, marginRight: 16 }}
        />
        <VariatesSelect
          variates={_.filter(variates, (o) => _.indexOf(o.scope, role) >= 0)}
        />
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
