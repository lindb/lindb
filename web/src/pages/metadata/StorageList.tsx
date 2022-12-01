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
  IconRedoStroked,
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
  Notification,
} from "@douyinfe/semi-ui";
import { StatusTip, StorageStatusView } from "@src/components";
import { Route } from "@src/constants";
import { Storage } from "@src/models";
import { ExecService } from "@src/services";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import { useQuery } from "@tanstack/react-query";
import React, { useContext } from "react";
import { UIContext } from "@src/context/UIContextProvider";

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
  const { locale } = useContext(UIContext);
  const { MetadataStorageView, Common } = locale;

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
  const recover = (storageName: string) => {
    return new Promise(async (resolve: any) => {
      try {
        const msg: any = await ExecService.exec({
          sql: `recover storage '${storageName}'`,
        });
        Notification.success({
          title: MetadataStorageView.recoverSuccessTitle,
          content: _.join(msg || [], ","),
          position: "top",
          theme: "light",
          duration: 5,
        });
        resolve();
      } catch (err) {
        Notification.error({
          title: MetadataStorageView.recoverErrorTitle,
          content: _.get(err, "response.data", Common.unknownInternalError),
          position: "top",
          theme: "light",
          duration: 5,
        });
        resolve();
      }
    });
  };
  const columns = [
    {
      title: MetadataStorageView.name,
      dataIndex: "config.namespace",
      width: 170,
      key: "name",
    },
    {
      title: MetadataStorageView.status,
      dataIndex: "status",
      key: "status",
      width: 120,
      render: (item: any) => {
        return <StorageStatusView text={item} showBadge={true} />;
      },
    },
    {
      title: MetadataStorageView.configuration,
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
      title: Common.actions,
      key: "actions",
      width: 100,
      render: (_text: any, r: any) => {
        return (
          <Popconfirm
            title={Common.pleaseConfirm}
            content={MetadataStorageView.recoverConfirmMessage}
            onConfirm={() => recover(_.get(r, "config.namespace"))}
          >
            <Button icon={<IconRedoStroked />} type="primary" />
          </Popconfirm>
        );
      },
    },
  ];

  return (
    <Card>
      <SplitButtonGroup style={{ marginBottom: 20 }}>
        <RegisterBtn text={MetadataStorageView.register} />
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
