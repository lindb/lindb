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
  const { isDark, collapsed, toggleCollapse } = useContext(UIContext);

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

  return (
    <Sider
    // conflict local setting
    // breakpoint={["lg"]}
    // onBreakpoint={(_screen, bool) => {
    //   UIStore.setSidebarCollapse(!bool);
    // }}
    >
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
        items={menus as any[]}
        selectedKeys={selectedKeys}
        onClick={(data) => {
          const item = routes.get(`${data.itemKey}`);
          const needClearKeys = _.pullAll(
            URLStore.getParamKeys(),
            _.get(item, "keep", [])
          );
          URLStore.changeURLParams({
            path: item?.path,
            needDelete: needClearKeys,
          });
        }}
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
      />
    </Sider>
  );
};

export default SiderMenu;
