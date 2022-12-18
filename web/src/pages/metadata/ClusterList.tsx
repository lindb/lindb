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
import { IconPlusCircle, IconRefresh } from "@douyinfe/semi-icons";
import {
  Button,
  Card,
  Descriptions,
  SplitButtonGroup,
  Table,
} from "@douyinfe/semi-ui";
import { ClusterStatusView, StatusTip } from "@src/components";
import { Storage } from "@src/models";
import { ExecService } from "@src/services";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import { useQuery } from "@tanstack/react-query";
import React, { useContext } from "react";
import { UIContext } from "@src/context/UIContextProvider";

export const ClusterList: React.FC<{
  sql: string;
  registerPath: string;
  actions?: any;
}> = (props) => {
  const { sql, registerPath, actions } = props;
  const { isLoading, isFetching, data, refetch, error, isError } = useQuery(
    ["show_cluster"],
    async () => {
      return ExecService.exec<Storage[]>({ sql: sql });
    }
  );
  const { locale } = useContext(UIContext);
  const { MetadataClusterView } = locale;

  const RegisterBtn: React.FC<any> = ({ text }) => {
    return (
      <Button
        icon={<IconPlusCircle />}
        onClick={() => {
          URLStore.changeURLParams({ path: registerPath });
        }}
      >
        {text}
      </Button>
    );
  };
  const columns = [
    {
      title: MetadataClusterView.name,
      dataIndex: "config.namespace",
      width: 170,
      key: "name",
    },
    {
      title: MetadataClusterView.status,
      dataIndex: "status",
      key: "status",
      width: 120,
      render: (item: any) => {
        return <ClusterStatusView text={item} showBadge={true} />;
      },
    },
    {
      title: MetadataClusterView.configuration,
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
  ];

  return (
    <Card>
      <SplitButtonGroup style={{ marginBottom: 20 }}>
        <RegisterBtn text={MetadataClusterView.register} />
        <Button
          icon={<IconRefresh />}
          style={{ marginLeft: 4 }}
          onClick={() => refetch()}
        />
      </SplitButtonGroup>
      <Table
        className="lin-table"
        dataSource={data}
        loading={isLoading || isFetching}
        empty={
          <StatusTip
            isLoading={isLoading || isFetching}
            isError={isError}
            isEmpty={_.isEmpty(data)}
            error={error}
          />
        }
        columns={actions ? [...columns, actions] : columns}
        pagination={false}
      />
    </Card>
  );
};
