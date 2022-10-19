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
import { Card, Descriptions, Typography } from "@douyinfe/semi-ui";
import { DateTimeFormat } from "@src/constants";
import { Master } from "@src/models";
import { ExecService } from "@src/services";
import { useQuery } from "@tanstack/react-query";
import moment from "moment";
import React from "react";
import { StatusTip } from "@src/components";

const { Text } = Typography;

const MasterView: React.FC = () => {
  const {
    isLoading,
    data: master,
    isError,
    error,
  } = useQuery(["show_master"], async () => {
    return ExecService.exec<Master>({ sql: "show master" });
  });

  const renderMaster = () => {
    if (isLoading || isError) {
      return (
        <StatusTip isLoading={isLoading} isError={isError} error={error} />
      );
    }

    return (
      <Descriptions
        className="lin-description"
        row
        size="small"
        data={[
          {
            key: "Elect Time",
            value: (
              <Text link>
                {master?.electTime &&
                  moment(master?.electTime).format(DateTimeFormat)}
              </Text>
            ),
          },
          { key: "Host IP", value: <Text link>{master?.node?.hostIp}</Text> },
          {
            key: "Host Name",
            value: <Text link>{master?.node?.hostName}</Text>,
          },
          {
            key: "GRPC Port",
            value: <Text link>{master?.node?.grpcPort}</Text>,
          },
          {
            key: "HTTP Port",
            value: <Text link>{master?.node?.httpPort}</Text>,
          },
        ]}
      />
    );
  };

  return (
    <>
      <Card
        title="Master"
        headerStyle={{ padding: 12 }}
        bodyStyle={{ padding: 12 }}
      >
        {renderMaster()}
      </Card>
    </>
  );
};

export default MasterView;
