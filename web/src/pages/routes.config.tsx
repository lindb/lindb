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
  IconComponentPlaceholderStroked,
  IconElementStroked,
  IconGridStroked,
  IconFixedStroked,
  IconChecklistStroked,
  IconSearchStroked,
  IconSendStroked,
  IconShareStroked,
  IconSettingStroked,
  IconMonitorStroked,
} from "@douyinfe/semi-icons";
import { CommonVariates, Dashboards } from "@src/configs";
import { Route } from "@src/constants";
import { Icon } from "@src/components";
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
  RequestView,
  BrokerList,
  BrokerConfig,
  LogicDatabaseList,
  LogicDatabaseConfig,
} from "@src/pages";
import * as _ from "lodash-es";
import React from "react";
import { RouteItem } from "@src/models";

export const routes = [
  {
    text: "Overview",
    path: Route.Overview,
    icon: <IconComponentPlaceholderStroked size="large" />,
    content: <Overview />,
    help: "/guide/admin-ui/overview.html",
    items: [
      {
        inner: true,
        itemKey: "Overview/Storage",
        text: "Storage",
        path: Route.StorageOverview,
        content: <StorageOverview />,
        help: "/guide/admin-ui/overview.html#storage-cluster-status",
      },
      {
        inner: true,
        itemKey: "Overview/Configuration",
        text: "Configuration",
        path: Route.ConfigurationView,
        content: <ConfigurationView />,
        help: "/guide/admin-ui/overview.html#node-configuration",
      },
    ],
  },
  {
    text: "Search",
    path: Route.Search,
    icon: <IconSearchStroked size="large" />,
    content: <DataSearch />,
    help: "/guide/admin-ui/search.html",
  },
  {
    text: "Explore",
    path: Route.Explore,
    timePicker: true,
    icon: <Icon icon="iconExplore" style={{ fontSize: 20 }} />,
    content: <DataExplore />,
    help: "/guide/admin-ui/explore.html",
  },
  {
    text: "Monitoring",
    itemKey: "Monitoring",
    icon: <IconMonitorStroked size="large" />,
    items: [
      {
        text: "Dashboard",
        path: Route.MonitoringDashboard,
        icon: <IconGridStroked size="large" />,
        content: (
          <DashboardView variates={CommonVariates} dashboards={Dashboards} />
        ),
        timePicker: true,
        help: "/guide/admin-ui/monitoring.html#dashboard",
      },
      {
        text: "Replication",
        roles: ["Broker"],
        path: Route.MonitoringReplication,
        icon: <IconElementStroked size="large" />,
        content: <ReplicationView />,
        help: "/guide/admin-ui/monitoring.html#replication",
      },
      {
        text: "Request",
        path: Route.MonitoringRequest,
        icon: <IconSendStroked size="large" />,
        content: <RequestView />,
        help: "/guide/admin-ui/monitoring.html#request",
      },
      {
        text: "Log View",
        path: Route.MonitoringLogs,
        icon: <IconChecklistStroked size="large" />,
        content: <LogView />,
        help: "/guide/admin-ui/monitoring.html#log-view",
      },
    ],
  },
  {
    text: "Metadata",
    itemKey: "Metadata",
    icon: <IconSettingStroked size="large" />,
    items: [
      {
        text: "Storage",
        path: Route.MetadataStorage,
        roles: ["Broker"],
        icon: <Icon icon="iconts-tubiao_APPCluster" style={{ fontSize: 20 }} />,
        content: <StorageList />,
        help: "/guide/admin-ui/metadata.html#storage",
        items: [
          {
            inner: true,
            itemKey: "Metadata/Storage/Configuration",
            text: "Configuration",
            path: Route.MetadataStorageConfig,
            content: <StorageConfig />,
            help: "/guide/admin-ui/metadata.html#storage",
          },
        ],
      },
      {
        text: "Broker",
        path: Route.MetadataBroker,
        roles: ["Root"],
        icon: <Icon icon="iconts-tubiao_APPCluster" style={{ fontSize: 20 }} />,
        content: <BrokerList />,
        help: "/guide/admin-ui/metadata.html#storage",
        items: [
          {
            roles: ["Root"],
            inner: true,
            itemKey: "Metadata/Broker/Configuration",
            text: "Configuration",
            path: Route.MetadataBrokerConfig,
            content: <BrokerConfig />,
            help: "/guide/admin-ui/metadata.html#storage",
          },
        ],
      },
      {
        text: "Database",
        path: Route.MetadataDatabase,
        roles: ["Broker"],
        icon: <Icon icon="icondatabase" style={{ fontSize: 20 }} />,
        content: <DatabaseList />,
        help: "/guide/admin-ui/metadata.html#database",
        items: [
          {
            inner: true,
            itemKey: "Metadata/Database/Configuration",
            text: "Configuration",
            path: Route.MetadataDatabaseConfig,
            content: <DatabaseConfig />,
            help: "/guide/admin-ui/metadata.html#database",
          },
        ],
      },
      {
        text: "Logic Database",
        roles: ["Root"],
        path: Route.MetadataLogicDatabase,
        icon: <Icon icon="icondatabase" style={{ fontSize: 20 }} />,
        content: <LogicDatabaseList />,
        help: "/guide/admin-ui/metadata.html#database",
        items: [
          {
            inner: true,
            roles: ["Root"],
            itemKey: "Metadata/Logic/Database/Configuration",
            text: "Configuration",
            path: Route.MetadataLogicDatabaseConfig,
            content: <LogicDatabaseConfig />,
            help: "/guide/admin-ui/metadata.html#database",
          },
        ],
      },
      {
        text: "Explore",
        path: Route.MetadataExplore,
        icon: <IconFixedStroked size="large" />,
        content: <MetadataExplore />,
        help: "/guide/admin-ui/metadata.html#explore",
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
        rs.push({ ...item, itemKey: item.itemKey || item.path });
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
