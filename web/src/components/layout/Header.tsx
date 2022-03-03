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
  IconGithubLogo,
  IconHelpCircleStroked,
  IconHomeStroked,
} from "@douyinfe/semi-icons";
import { Breadcrumb, Button, Layout, Nav, Space } from "@douyinfe/semi-ui";
import { RouteItem, routeMap } from "@src/configs";
import React, { useEffect, useState } from "react";
import { TimePicker } from "@src/components";
import * as _ from "lodash-es";
import { URLStore } from "@src/stores";
import { useWatchURLChange } from "@src/hooks";
const { Header: HeaderUI } = Layout;

export default function Header() {
  const [breadcrumbRoutes, setBreadcrumbRoutes] = useState<any[]>([]);
  const [currentRouter, setCurrentRouter] = useState<any>({});

  useWatchURLChange(() => {
    const pathname = URLStore.path;
    const currentRouter = routeMap.get(pathname || "");
    const breadcrumbItems: any[] = [];
    if (currentRouter) {
      const generate = (item: RouteItem) => {
        if (item.parnet) {
          generate(item.parnet);
        }
        breadcrumbItems.push({
          href: item.content ? `#${item.path}` : null,
          name: item.text,
          icon: !item.parnet ? <IconHomeStroked /> : null,
        });
      };

      generate(currentRouter);
    }
    setCurrentRouter(currentRouter);
    setBreadcrumbRoutes(breadcrumbItems);
  });

  return (
    <HeaderUI
      style={{
        position: "fixed",
        width: "calc(100% - 220px)",
        left: 220,
        top: 0,
        right: 0,
        zIndex: 1000,
      }}
    >
      <Nav
        mode="horizontal"
        style={{ paddingRight: 12, paddingLeft: 16 }}
        header={
          <Space align="center">
            <Breadcrumb
              className="lin-header-breadcrumb"
              routes={breadcrumbRoutes}
            />
          </Space>
        }
        footer={
          <>
            {_.get(currentRouter, "timePicker", false) && <TimePicker />}
            <Button
              icon={<IconGithubLogo size="large" />}
              style={{
                color: "var(--semi-color-text-2)",
                marginRight: 12,
                marginLeft: 12,
              }}
              onClick={() => window.open("https://github.com/lindb/lindb")}
            />
            <Button
              icon={<IconHelpCircleStroked size="large" />}
              style={{
                color: "var(--semi-color-text-2)",
              }}
              onClick={() => window.open("https://lindb.io")}
            />
          </>
        }
      ></Nav>
    </HeaderUI>
  );
}
