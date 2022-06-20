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
  IconAppCenter,
  IconComponentPlaceholderStroked,
  IconElementStroked,
  IconGridStroked,
  IconInherit,
  IconListView,
  IconSearch,
  IconServer,
  IconTemplate,
  IconShareStroked,
} from "@douyinfe/semi-icons";
import { CommonVariates, Dashboards } from "@src/configs";
import { Route } from "@src/constants";
import {
  Overview,
  ConfigurationView,
  DataExplore,
  DataSearch,
  ReplicationView,
  DashboardView,
  LogView,
  DatabaseConfig,
  DatabaseList,
  StorageConfig,
  StorageList,
  StorageOverview,
  MetadataExplore,
  MultipleIDCList,
} from "@src/pages";
import * as _ from "lodash-es";
import React from "react";

export type RouteItem = {
  itemKey?: string;
  parent?: RouteItem | null;
  text: string;
  path?: string;
  inner?: boolean;
  icon?: React.ReactNode | string;
  content?: React.FunctionComponent;
  items?: RouteItem[];
  keep?: string[]; //keep url params when change route
};

export const routes = [
  {
    text: "Overview",
    path: Route.Overview,
    icon: <IconComponentPlaceholderStroked size="large" />,
    content: <Overview />,
    items: [
      {
        inner: true,
        itemKey: "Overview/Storage",
        text: "Storage",
        path: Route.StorageOverview,
        content: <StorageOverview />,
      },
      {
        inner: true,
        itemKey: "Overview/Configuration",
        text: "Configuration",
        path: Route.ConfigurationView,
        content: <ConfigurationView />,
      },
    ],
  },
  {
    text: "Search",
    path: Route.Search,
    icon: <IconSearch size="large" />,
    content: <DataSearch />,
  },
  {
    text: "Explore",
    path: Route.Explore,
    timePicker: true,
    icon: <IconAppCenter size="large" />,
    content: <DataExplore />,
  },
  {
    text: "Monitoring",
    items: [
      {
        text: "Dashboard",
        path: Route.MonitoringDashboard,
        icon: <IconGridStroked size="large" />,
        content: (
          <DashboardView variates={CommonVariates} dashboards={Dashboards} />
        ),
        timePicker: true,
        keep: ["start", "end", "node", "db"],
      },
      {
        text: "Replication",
        path: Route.MonitoringReplication,
        icon: <IconElementStroked size="large" />,
        content: <ReplicationView />,
      },
      {
        text: "Log View",
        path: Route.MonitoringLogs,
        icon: <IconListView size="large" />,
        content: <LogView />,
      },
    ],
  },
  {
    text: "Metadata",
    items: [
      {
        text: "Storage",
        path: Route.MetadataStorage,
        icon: <IconTemplate size="large" />,
        content: <StorageList />,
        items: [
          {
            inner: true,
            itemKey: "Metadata/Storage/Configuration",
            text: "Configuration",
            path: Route.MetadataStorageConfig,
            content: <StorageConfig />,
          },
        ],
      },
      {
        text: "Database",
        path: Route.MetadataDatabase,
        icon: <IconServer size="large" />,
        content: <DatabaseList />,
        items: [
          {
            inner: true,
            itemKey: "Metadata/Database/Configuration",
            text: "Configuration",
            path: Route.MetadataDatabaseConfig,
            content: <DatabaseConfig />,
          },
        ],
      },
      {
        text: "Explore",
        path: Route.MetadataExplore,
        icon: <IconInherit size="large" />,
        content: <MetadataExplore />,
      },
      {
        text: "Multiple IDCs",
        path: Route.MetadataMultipleIDC,
        icon: <IconShareStroked size="large" />,
        content: <MultipleIDCList />,
      },
    ],
  },
] as unknown as RouteItem[];

function flattenRouters(routeItems: RouteItem[]): Map<string, RouteItem> {
  const rs = new Map<string, RouteItem>();
  const flatten = (items: RouteItem[], parent: RouteItem | null) => {
    items.map((item: RouteItem) => {
      item.parent = parent;
      if (item.items) {
        flatten(item.items, item);
      }
      if (item.path) {
        rs.set(item.path, item);
      }
    });
  };
  flatten(routeItems, null);
  return rs;
}

function getSwithRouterList(routeItems: RouteItem[]): RouteItem[] {
  const rs: RouteItem[] = [];
  const flatten = (items: RouteItem[]) => {
    items.map((item: RouteItem) => {
      if (item.items) {
        flatten(item.items);
      }
      if (item.content) {
        rs.push(item);
      }
    });
  };
  flatten(routeItems);
  return rs;
}

function getMenuList(routeItems: RouteItem[]): RouteItem[] {
  const forEach = (items: RouteItem[]): RouteItem[] => {
    const rs: RouteItem[] = [];
    items.map((item: RouteItem) => {
      if (item.items) {
        item.items = forEach(item.items);
      }
      if (!item.inner) {
        rs.push({ ...item, itemKey: item.path });
      }
    });
    return rs;
  };
  return forEach(routeItems);
}

function getDefaultOpenKeys(menus: any[]): string[] {
  return (menus || []).reduce((pre, item) => {
    if (item.items) {
      pre.push(item.itemKey);
      const newArray: string[] = pre.concat(
        getDefaultOpenKeys(item.items) || []
      );
      return newArray;
    }
    return pre;
  }, [] as string[]);
}

export const menus = getMenuList(_.cloneDeep(routes));
export const switchRouters = getSwithRouterList(_.cloneDeep(routes));
export const defaultOpenKeys = getDefaultOpenKeys(_.cloneDeep(routes));
export const routeMap = flattenRouters(_.cloneDeep(routes));
