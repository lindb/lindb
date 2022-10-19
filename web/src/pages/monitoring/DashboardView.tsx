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
import { Card, Form, Col, Row } from "@douyinfe/semi-ui";
import { Chart, LinSelect, MetadataSelect } from "@src/components";
import { StateRoleName } from "@src/constants";
import { useParams } from "@src/hooks";
import { ChartType, Dashboard, DashboardItem } from "@src/models";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { useMemo } from "react";

const DashboardForm: React.FC<{
  dashboards?: DashboardItem[];
  variates?: any;
}> = (props) => {
  const { dashboards, variates } = props;
  const { role } = useParams(["role"]);

  const dashboardForRole = _.filter(
    dashboards,
    (o) => _.indexOf(o.scope, role || StateRoleName.Broker) >= 0
  );
  return (
    <Form
      className="lin-variate-form"
      style={{ paddingBottom: 0, paddingTop: 0, display: "inline-flex" }}
      layout="horizontal"
    >
      <LinSelect
        showClear
        field="role"
        defaultValue={StateRoleName.Broker}
        prefix={<IconFixedStroked />}
        loader={() => [
          { value: StateRoleName.Broker, label: StateRoleName.Broker },
          { value: StateRoleName.Storage, label: StateRoleName.Storage },
        ]}
        clearKeys={["namespace", "node", "d"]}
        style={{ minWidth: 60 }}
      />
      <LinSelect
        field="d"
        defaultValue={_.get(dashboardForRole, "[0].value")}
        prefix={<IconGridStroked />}
        loader={() => {
          const params = URLStore.getParams();
          return _.filter(
            dashboards,
            (o) =>
              _.indexOf(o.scope, _.get(params, "role", StateRoleName.Broker)) >=
              0
          );
        }}
        style={{ minWidth: 60 }}
        reloadKeys={["role"]}
      />
      {_.map(
        _.filter(
          variates,
          (o) => _.indexOf(o.scope, role || StateRoleName.Broker) >= 0
        ),
        (v: any) => (
          <MetadataSelect
            key={v.tagKey}
            variate={v}
            multiple={v.multiple}
            type="tagValue"
          />
        )
      )}
    </Form>
  );
};

const ViewDashboard: React.FC<{
  dashboards: DashboardItem[];
}> = (props) => {
  const { dashboards } = props;
  const { d, role } = useParams(["d", "role"]);

  const dashboard = useMemo(() => {
    const r = role || StateRoleName.Broker;
    const currentDashboards = _.filter(
      dashboards,
      (o) => _.indexOf(o.scope, r) >= 0
    );
    let dashboardItem = _.find(currentDashboards, (item: DashboardItem) => {
      return item.value == d;
    });
    if (_.isUndefined(dashboardItem) || _.isNull(dashboardItem)) {
      dashboardItem = currentDashboards[0];
    }
    return _.get(dashboardItem, "dashboard") as Dashboard;
  }, [d, role, dashboards]);

  return (
    <>
      {(dashboard.rows || []).map((row: any, rowIdx: any) => (
        <Row
          style={{ marginTop: 12 }}
          key={rowIdx}
          gutter={dashboard.gutter || 8}
        >
          {row?.panels.map((panel: any, panelIdx: any) => (
            <Col span={panel.span || 12} key={`${rowIdx}-${panelIdx}`}>
              <Chart
                type={_.get(panel.chart, "config.type", ChartType.Line)}
                queries={panel.chart.targets || []}
                config={panel.chart}
              />
            </Col>
          ))}
        </Row>
      ))}
    </>
  );
};

const DashboardView: React.FC<{
  dashboards?: DashboardItem[];
  variates?: any;
}> = (props) => {
  const { dashboards } = props;

  return (
    <>
      <Card bodyStyle={{ padding: 12 }}>
        <DashboardForm {...props} />
      </Card>
      <ViewDashboard dashboards={dashboards || []} />
    </>
  );
};

export default DashboardView;
