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
  Notification,
} from "@douyinfe/semi-ui";
import { Route } from "@src/constants";
import { Database } from "@src/models";
import { exec } from "@src/services";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { useCallback, useEffect, useState } from "react";

const { Text } = Typography;

export default function MultipleIDCList() {
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

  const dropDatabase = useCallback(async (name) => {
    setError("");
    try {
      const rs = await exec<string>({ sql: `drop database '${name}'` });
      Notification.success({
        content: `${rs}`,
        position: "top",
        duration: 5,
        theme: "light",
      });
      getDatabaseList();
    } catch (err) {
      Notification.error({
        title: "Drop database error",
        content: _.get(err, "response.data", "Unknown internal error"),
        position: "top",
        theme: "light",
        duration: 5,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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

  return <div>come soon</div>;
}
