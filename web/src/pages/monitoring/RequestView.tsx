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
import React, { useContext } from "react";
import { Route, SQL } from "@src/constants";
import { ExecService } from "@src/services";
import { Request, Unit } from "@src/models";
import {
  Card,
  Table,
  SplitButtonGroup,
  Button,
  Tooltip,
} from "@douyinfe/semi-ui";
import * as _ from "lodash-es";
import { IconRefresh, IconPlay } from "@douyinfe/semi-icons";
import moment from "moment";
import { useQuery } from "@tanstack/react-query";
import { StatusTip } from "@src/components";
import { FormatKit } from "@src/utils";
import { UIContext } from "@src/context/UIContextProvider";

const RequestView: React.FC = () => {
  const {
    isLoading,
    isError,
    error,
    data: requests,
    refetch,
  } = useQuery(["show_requests"], async () => {
    return ExecService.exec<Request[]>({
      sql: SQL.ShowRequests,
    });
  });
  const { locale } = useContext(UIContext);
  const { RequestView, Common } = locale;

  const columns: any[] = [
    {
      title: RequestView.timestamp,
      dataIndex: "start",
      render: (start: number) => {
        const dateTime = moment(parseInt(`${start / 1000000}`));
        return dateTime.format("YYYY-MM-DD HH:mm:ss");
      },
    },
    {
      title: RequestView.duration,
      dataIndex: "duration",
      render: (_text: any, record: Request, _index: any) => {
        return FormatKit.format(
          new Date().getTime() - record.start / 1000000,
          Unit.Milliseconds
        );
      },
    },
    {
      title: RequestView.linQL,
      dataIndex: "sql",
      key: "lql",
    },
    {
      title: RequestView.database,
      dataIndex: "db",
    },
    {
      title: RequestView.entry,
      dataIndex: "entry",
    },
    {
      title: Common.actions,
      key: "actions",
      render: (_text: any, record: Request, _index: any) => {
        return (
          <Tooltip content={RequestView.runLinQL}>
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

  return (
    <Card>
      <SplitButtonGroup style={{ marginBottom: 20 }}>
        <Button icon={<IconRefresh />} onClick={() => refetch()} />
      </SplitButtonGroup>
      <Table
        className="lin-table"
        dataSource={requests || []}
        pagination={false}
        columns={columns}
        rowKey="requestId"
        empty={
          <StatusTip
            style={{ marginTop: 32 }}
            isLoading={isLoading}
            isError={isError}
            error={error}
          />
        }
      />
    </Card>
  );
};

export default RequestView;
