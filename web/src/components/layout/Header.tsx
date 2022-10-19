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
import React, { useState, useContext } from "react";
import {
  IconGithubLogo,
  IconHelpCircleStroked,
  IconHomeStroked,
} from "@douyinfe/semi-icons";
import { Breadcrumb, Button, Layout, Nav, Space } from "@douyinfe/semi-ui";
import { TimePicker, Icon } from "@src/components";
import { RouteItem } from "@src/models";
import { useWatchURLChange } from "@src/hooks";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import { UIContext } from "@src/context/UIContextProvider";
const { Header: HeaderUI } = Layout;

const Header: React.FC<{ routes: Map<string, RouteItem> }> = (props) => {
  const { routes } = props;
  const { toggleTheme, collapsed, isDark } = useContext(UIContext);
  const [breadcrumbRoutes, setBreadcrumbRoutes] = useState<any[]>([]);
  const [currentRouter, setCurrentRouter] = useState<any>({});

  useWatchURLChange(() => {
    const pathname = URLStore.getPath();
    const currentRouter = routes.get(pathname || "");
    const breadcrumbItems: any[] = [];
    if (currentRouter) {
      const generate = (item: RouteItem) => {
        if (item.parent) {
          generate(item.parent);
        }
        breadcrumbItems.push({
          href: item.content ? `#${item.path}` : null,
          name: item.text,
          icon: !item.parent ? <IconHomeStroked /> : null,
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
        left: collapsed ? 60 : 220,
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
              icon={
                isDark() ? (
                  <Icon icon="iconmoon-line" />
                ) : (
                  <Icon icon="iconSun" />
                )
              }
              style={{
                color: "var(--semi-color-text-2)",
                marginLeft: 12,
              }}
              onClick={() => {
                toggleTheme();
              }}
            />
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
};

export default Header;
