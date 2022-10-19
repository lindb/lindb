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
import { Layout } from "@douyinfe/semi-ui";
import { Footer, Header, SiderMenu } from "@src/components";
import { defaultOpenKeys, menus, routeMap, switchRouters } from "@src/pages";
import React from "react";
import { Redirect, Route, Switch } from "react-router-dom";

const { Content } = Layout;

export default function Console() {
  return (
    <Layout style={{ height: "100vh" }}>
      <SiderMenu
        defaultOpenAll
        openKeys={defaultOpenKeys}
        routes={routeMap}
        menus={menus}
      />
      <Layout>
        <Header routes={routeMap} />
        <Content
          style={{
            padding: "71px 12px 12px",
            backgroundColor: "var(--semi-color-bg-0)",
          }}
        >
          <Switch>
            {switchRouters.map((item) => (
              <Route
                key={item.path}
                path={item.path}
                render={() => item.content}
              />
            ))}
            <Redirect to="/overview" />
          </Switch>
        </Content>
        <Footer />
      </Layout>
    </Layout>
  );
}
