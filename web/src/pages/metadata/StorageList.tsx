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
  IconDeleteStroked,
  IconPlusCircle,
  IconRefresh,
} from "@douyinfe/semi-icons";
import {
  Button,
  Card,
  Descriptions,
  Popconfirm,
  SplitButtonGroup,
  Table,
} from "@douyinfe/semi-ui";
import { StatusTip, StorageStatusView } from "@src/components";
import { Route } from "@src/constants";
import { Storage } from "@src/models";
import { ExecService } from "@src/services";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import { useQuery } from "@tanstack/react-query";
import React from "react";

const columns = [
  {
    title: "Name(Namespace)",
    dataIndex: "config.namespace",
    width: 170,
    key: "name",
  },
  {
    title: "Status",
    dataIndex: "status",
    key: "status",
    width: 100,
    render: (item: any) => {
      return <StorageStatusView text={item} showBadge={true} />;
    },
  },
  {
    title: "Configuration",
    dataIndex: "config",
    key: "config",
    render: (item: any) => {
      var configItems: any[] = [];
      Object.keys(item).forEach(function (key) {
        const val = item[key];
        if (val) {
          configItems.push({
            key: `${key}:`,
            value: Array.isArray(val) ? JSON.stringify(val) : val,
          });
        }
      });
      return <Descriptions data={configItems} size="small" />;
    },
  },
  {
    title: "Actions",
    key: "actions",
    width: 100,
    render: () => {
      return (
        <Popconfirm
          title="Please confirm"
          content="Are you sure want to remove storage?"
        >
          <Button icon={<IconDeleteStroked />} type="danger" />
        </Popconfirm>
      );
    },
  },
];
export default function StorageList() {
  const {
    isLoading,
    isFetching,
    data: storageList,
    refetch,
    error,
    isError,
  } = useQuery(["show_storage"], async () => {
    return ExecService.exec<Storage[]>({ sql: "show storages" });
  });

  const RegisterBtn: React.FC<any> = ({ text }) => {
    return (
      <Button
        icon={<IconPlusCircle />}
        onClick={() => {
          URLStore.changeURLParams({ path: Route.MetadataStorageConfig });
        }}
      >
        {text}
      </Button>
    );
  };

  return (
    <Card>
      <SplitButtonGroup style={{ marginBottom: 20 }}>
        <RegisterBtn text="Register" />
        <Button
          icon={<IconRefresh />}
          style={{ marginLeft: 4 }}
          onClick={() => refetch()}
        />
      </SplitButtonGroup>
      <Table
        className="lin-table"
        dataSource={storageList}
        loading={isLoading || isFetching}
        empty={
          <StatusTip
            isLoading={isLoading || isFetching}
            isError={isError}
            isEmpty={_.isEmpty(storageList)}
            error={error}
          />
        }
        columns={columns}
        pagination={false}
      />
    </Card>
  );
}
