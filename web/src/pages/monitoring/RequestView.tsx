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
import React, { useState, useEffect } from "react";
import { Route, SQL } from "@src/constants";
import { exec } from "@src/services";
import { Request } from "@src/models";
import {
  Card,
  Table,
  SplitButtonGroup,
  Button,
  Empty,
  Typography,
  Tooltip,
} from "@douyinfe/semi-ui";
import * as _ from "lodash-es";
import { IconRefresh, IconPlay } from "@douyinfe/semi-icons";
import moment from "moment";
import {
  IllustrationConstructionDark,
  IllustrationNoContentDark,
} from "@douyinfe/semi-illustrations";

const Text = Typography.Text;

export default function RequestView() {
  const [requests, setRequests] = useState<Request[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const columns: any[] = [
    {
      title: "Timestamp",
      dataIndex: "start",
      render: (start: number) => {
        const dateTime = moment(parseInt(`${start / 1000000}`));
        return dateTime.format("YYYY-MM-DD HH:mm:ss");
      },
    },
    {
      title: "LQL",
      dataIndex: "sql",
      key: "lql",
    },
    {
      title: "Database",
      dataIndex: "db",
    },
    {
      title: "Broker",
      dataIndex: "broker",
    },
    {
      title: "Actions",
      key: "actions",
      render: (_text: any, record: Request, _index: any) => {
        return (
          <Tooltip content="Run Lin Query Language">
            <Button
              icon={<IconPlay />}
              style={{ color: "var(--semi-color-success)" }}
              onClick={() => {
                const url = _.split(window.location.href, "#", 1)[0];
                const sql = encodeURIComponent(record.sql);
                const path = `#${Route.Search}?db=${record.db}&sql=${sql}`;
                window.open(url + path, "_blank");
              }}
            />
          </Tooltip>
        );
      },
    },
  ];

  // get request list
  const fetchRequests = async () => {
    setLoading(true);
    try {
      const requests = await exec<Request[]>({ sql: SQL.ShowRequests });
      setRequests(requests);
    } catch (err) {
      setError(_.get(err, "response.data", "Unknown internal error"));
      setRequests([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRequests();
  }, []);

  return (
    <Card>
      <SplitButtonGroup style={{ marginBottom: 20 }}>
        <Button icon={<IconRefresh />} onClick={fetchRequests} />
      </SplitButtonGroup>
      <Table
        className="lin-table"
        dataSource={requests}
        loading={loading}
        pagination={false}
        columns={columns}
        rowKey="requestId"
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
                  <Text link onClick={fetchRequests}>
                    retry
                  </Text>{" "}
                  later.
                </Text>
              ) : (
                "No alive request"
              )
            }
          ></Empty>
        }
      />
    </Card>
  );
}
