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
  Badge,
  Button,
  Card,
  Descriptions,
  Empty,
  Space,
  SplitButtonGroup,
  Table,
  Typography,
  Popconfirm,
} from "@douyinfe/semi-ui";
import { Route } from "@src/constants";
import { Storage } from "@src/models";
import { exec } from "@src/services";
import React, { useCallback, useEffect, useState } from "react";
import { URLStore } from "@src/stores";

const { Text } = Typography;
const StorageStatus: React.FC<any> = ({ color, text }) => {
  return (
    <Space align="center">
      <Badge
        dot
        style={{
          width: 12,
          height: 12,
          marginTop: 4,
          backgroundColor: `var(--semi-color-${color})`,
        }}
      />
      <Text style={{ color: `var(--semi-color-${color})` }}> {text}</Text>
    </Space>
  );
};
const columns = [
  {
    title: "Name(Namespace)",
    dataIndex: "config.namespace",
    key: "name",
  },
  {
    title: "Status",
    dataIndex: "status",
    key: "status",
    render: (item: any) => {
      switch (item) {
        case "Ready":
          return <StorageStatus color="success" text={item} />;
        case "Initialize":
          return <StorageStatus color="secondary" text={item} />;
        default:
          return <StorageStatus color="warning" text={item} />;
      }
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
  const [storageList, setStorageList] = useState([] as Storage[]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const getStorageList = useCallback(async () => {
    setError("");
    setLoading(true);
    try {
      const list = await exec<Storage[]>({ sql: "show storages" });
      setStorageList(list || []);
    } catch (err) {
      setError(err?.message);
      setStorageList([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    getStorageList();
  }, [getStorageList]);

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
    <Card bordered={false}>
      <SplitButtonGroup style={{ marginBottom: 20 }}>
        <RegisterBtn text="Register" />
        <Button
          icon={<IconRefresh />}
          style={{ marginLeft: 4 }}
          onClick={getStorageList}
        />
      </SplitButtonGroup>
      <Table
        className="lin-table"
        dataSource={storageList}
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
                  <Text link onClick={getStorageList}>
                    retry
                  </Text>{" "}
                  later.
                </Text>
              ) : (
                "No storage cluster, please register one."
              )
            }
          >
            {!error && <RegisterBtn text="Register Now" />}
          </Empty>
        }
        columns={columns}
        pagination={false}
      />
    </Card>
  );
}
