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
  IllustrationConstructionDark,
  IllustrationNoContentDark,
} from "@douyinfe/semi-illustrations";
import {
  Button,
  Card,
  Empty,
  Popconfirm,
  SplitButtonGroup,
  Table,
  Typography,
} from "@douyinfe/semi-ui";
import { Route } from "@src/constants";
import { Database } from "@src/models";
import { exec } from "@src/services";
import { URLStore } from "@src/stores";
import React, { useCallback, useEffect, useState } from "react";

const { Text } = Typography;

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
    render: () => {
      return (
        <Popconfirm
          title="Please confirm"
          content="Are you sure want to remove database?"
        >
          <Button icon={<IconDeleteStroked />} type="danger" />
        </Popconfirm>
      );
    },
  },
];
export default function DatabaseList() {
  const [databaseList, setDatabaseList] = useState([] as Database[]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const getDatabaseList = useCallback(async () => {
    setError("");
    setLoading(true);
    try {
      const list = await exec<Database[]>({ sql: "show schemas" });
      setDatabaseList(list || []);
    } catch (err) {
      setError(err?.message);
      setDatabaseList([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    getDatabaseList();
  }, [getDatabaseList]);

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
    <Card bordered={false}>
      <SplitButtonGroup style={{ marginBottom: 20 }}>
        <CreateBtn text="Create" />
        <Button
          icon={<IconRefresh />}
          style={{ marginLeft: 4 }}
          onClick={getDatabaseList}
        />
      </SplitButtonGroup>
      <Table
        className="lin-table"
        dataSource={databaseList}
        loading={loading}
        empty={
          <Empty
            image={
              error ? (
                <IllustrationConstructionDark />
              ) : (
                <IllustrationNoContentDark />
              )
            }
            imageStyle={{
              height: 60,
            }}
            title={
              error ? (
                <Text type="danger">
                  {error}, please{" "}
                  <Text link onClick={getDatabaseList}>
                    retry
                  </Text>{" "}
                  later.
                </Text>
              ) : (
                "No database, please create one."
              )
            }
          >
            {!error && <CreateBtn text="Create Now" />}
          </Empty>
        }
        columns={columns}
        pagination={false}
      />
    </Card>
  );
}
