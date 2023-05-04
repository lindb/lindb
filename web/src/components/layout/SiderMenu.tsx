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
import { Layout, Nav, Space } from "@douyinfe/semi-ui";
import DarkLogo from "@src/assets/logo_dark.svg";
import Logo from "@src/assets/logo.svg";
import { useWatchURLChange } from "@src/hooks";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { useState, useContext } from "react";
import { UIContext } from "@src/context/UIContextProvider";
import { RouteItem } from "@src/models";
const { Sider } = Layout;

const SiderMenu: React.FC<{
  defaultOpenAll?: boolean;
  openKeys: string[];
  routes: Map<string, RouteItem>;
  menus: RouteItem[];
}> = (props) => {
  const { defaultOpenAll, routes, menus, openKeys } = props;
  const [selectedKeys, setSelectedKeys] = useState([] as string[]);
  const { isDark, collapsed, toggleCollapse, locale, env } =
    useContext(UIContext);
  const { SiderMenu } = locale;

  useWatchURLChange(() => {
    const path = URLStore.path;
    let key = "";
    const findSelectedKeys = (menus: any[]) => {
      (menus || []).map((item) => {
        if (path.includes(item.path)) {
          if (key.length < item.path.length) {
            key = item.path;
          }
        }
        if (item.items) {
          findSelectedKeys(item.items);
        }
      });
    };
    findSelectedKeys(menus);
    setSelectedKeys([key]);
  });

  const renderMenus = (menus: any, level: number) => {
    return (menus || []).map((item: any) => {
      // need match role
      if (item.roles && item.roles.indexOf(env.role) < 0) {
        return;
      }
      const subItems = _.filter(item.items, (o) => !_.get(o, "inner", false));
      if (_.size(subItems) > 0) {
        return (
          <Nav.Sub
            isOpen
            key={item.itemKey}
            itemKey={item.itemKey}
            icon={item.icon}
            text={SiderMenu[item.text]}
          >
            {renderMenus(subItems, level + 1)}
          </Nav.Sub>
        );
      }
      return (
        <Nav.Item
          level={level}
          key={item.itemKey}
          itemKey={item.itemKey}
          icon={item.icon}
          text={SiderMenu[item.text]}
          onClick={() => {
            const routeItem = routes.get(`${item.itemKey}`);
            const needClearKeys = _.pullAll(
              URLStore.getParamKeys(),
              _.get(item, "keep", [])
            );
            URLStore.changeURLParams({
              path: routeItem?.path,
              needDelete: needClearKeys,
            });
          }}
        />
      );
    });
  };

  return (
    <Sider>
      <Nav
        className="lin-nav"
        defaultOpenKeys={defaultOpenAll ? openKeys : []}
        subNavMotion={false}
        limitIndent={false}
        isCollapsed={collapsed}
        onCollapseChange={() => toggleCollapse()}
        style={{
          maxWidth: 220,
          height: "100%",
        }}
        selectedKeys={selectedKeys}
        header={{
          logo: (
            <img
              src={isDark() ? DarkLogo : Logo}
              style={{ width: 48, height: 48, marginRight: 8 }}
            />
          ),
          text: (
            <Space align="end">
              <div
                style={{
                  fontSize: 32,
                  height: 32,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                LinDB
              </div>
            </Space>
          ),
          style: { paddingTop: 12, paddingBottom: 12, paddingLeft: 2 },
        }}
        footer={{
          collapseButton: true,
        }}
      >
        {renderMenus(menus, 0)}
      </Nav>
    </Sider>
  );
};

export default SiderMenu;
