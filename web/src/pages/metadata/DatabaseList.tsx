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
  Popconfirm,
  SplitButtonGroup,
  Table,
  Typography,
  Notification,
} from "@douyinfe/semi-ui";
import StatusTip from "@src/components/common/StatusTip";
import { Route } from "@src/constants";
import { Database } from "@src/models";
import { ExecService } from "@src/services";
import { URLStore } from "@src/stores";
import { useQuery } from "@tanstack/react-query";
import * as _ from "lodash-es";
import React from "react";

const { Text } = Typography;

export default function DatabaseList() {
  const { isLoading, isError, isFetching, error, data, refetch } = useQuery(
    ["show_database_schemas"],
    async () => {
      return ExecService.exec<Database[]>({ sql: "show schemas" });
    }
  );

  const dropDatabase = async (name: string) => {
    try {
      const rs = await ExecService.exec<string>({
        sql: `drop database '${name}'`,
      });
      Notification.success({
        content: `${rs}`,
        position: "top",
        duration: 5,
        theme: "light",
      });
      refetch();
    } catch (err) {
      Notification.error({
        title: "Drop database error",
        content: _.get(err, "response.data", "Unknown internal error"),
        position: "top",
        theme: "light",
        duration: 5,
      });
    }
  };

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Storage",
      dataIndex: "storage",
    },
    {
      title: "Description",
      dataIndex: "desc",
    },
    {
      title: "Actions",
      key: "actions",
      width: 100,
      render: (_text: any, record: any, _index: any) => {
        return (
          <Popconfirm
            title="Please confirm"
            content={
              <>
                Are you sure drop [
                <Text strong type="danger">
                  {record.name}
                </Text>
                ] database?
              </>
            }
            onConfirm={() => {
              dropDatabase(record.name);
            }}
          >
            <Button icon={<IconDeleteStroked />} type="danger" />
          </Popconfirm>
        );
      },
    },
  ];

  const CreateBtn: React.FC<any> = ({ text }) => {
    return (
      <Button
        icon={<IconPlusCircle />}
        onClick={() => {
          URLStore.changeURLParams({ path: Route.MetadataDatabaseConfig });
        }}
      >
        {text}
      </Button>
    );
  };

  return (
    <Card>
      <SplitButtonGroup style={{ marginBottom: 20 }}>
        <CreateBtn text="Create" />
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
        columns={columns}
        pagination={false}
      />
    </Card>
  );
}
