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
import { SQL } from "@src/constants";
import { ExecService } from "@src/services";
import { useQuery } from "@tanstack/react-query";
import React, { useContext } from "react";
import { UIContext } from "@src/context/UIContextProvider";
import { Card, Descriptions, Table, Typography } from "@douyinfe/semi-ui";
import * as _ from "lodash-es";
import { BrokerState } from "@src/models";
import { StatusTip } from "@src/components";

const { Text } = Typography;

const showBrokers = SQL.ShowBrokerAliveNodes;
const BrokerView: React.FC = () => {
  const { isLoading, isError, error, data } = useQuery(
    ["show_brokers"],
    async () => {
      return ExecService.exec<any[]>({ sql: showBrokers });
    }
  );
  const { locale } = useContext(UIContext);
  const { BrokerView } = locale;

  const columns = [
    {
      title: BrokerView.name,
      dataIndex: "name",
      key: "name",
    },
    {
      title: BrokerView.nodeStatus,
      render: (_text: any, record: BrokerState, _index: any) => {
        return (
          <Descriptions
            row
            className="lin-small-desc"
            size="small"
            data={[
              {
                key: BrokerView.aliveNodes,
                value: (
                  <Text link>
                    {_.keys(_.get(record, "liveNodes", {})).length}
                  </Text>
                ),
              },
            ]}
          />
        );
      },
    },
  ];
  console.log("data...", data);
  return (
    <Card
      title={BrokerView.brokerClusterList}
      headerStyle={{ padding: 12 }}
      bodyStyle={{ padding: 12 }}
    >
      <Table
        size="small"
        bordered={false}
        columns={columns}
        dataSource={data || []}
        loading={isLoading}
        pagination={false}
        empty={
          <StatusTip isLoading={isLoading} isError={isError} error={error} />
        }
      />
    </Card>
  );
};

export default BrokerView;
