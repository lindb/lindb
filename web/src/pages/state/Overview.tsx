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
import { MasterView, NodeView, StatusTip, StorageView } from "@src/components";
import { StateMetricName, SQL } from "@src/constants";
import { UIContext } from "@src/context/UIContextProvider";
import { useStorage } from "@src/hooks";
import { ExecService } from "@src/services";
import { useQuery } from "@tanstack/react-query";
import React, { useContext } from "react";

// must define outside function component, if define in component maybe endless loop.
const brokerMetric = `show broker metric where metric in ('${StateMetricName.CPU}','${StateMetricName.Memory}')`;

const Overview: React.FC = () => {
  const {
    isLoading: storageLoading,
    isError: storageHasError,
    error: storageError,
    storages,
  } = useStorage();
  const {
    isLoading: nodeLoading,
    data: liveNodes,
    isError: nodeHasError,
    error: nodeError,
  } = useQuery(["show_broker_alive_nodes"], async () => {
    return ExecService.exec<any[]>({ sql: SQL.ShowBrokerAliveNodes });
  });
  const { locale } = useContext(UIContext);
  const { Overview } = locale;

  return (
    <>
      <MasterView />
      <NodeView
        title={Overview.brokerLiveNodes}
        nodes={liveNodes || []}
        sql={brokerMetric}
        style={{ marginTop: 12, marginBottom: 12 }}
        statusTip={
          <StatusTip
            isLoading={nodeLoading}
            isError={nodeHasError}
            error={nodeError}
          />
        }
      />
      <StorageView
        storages={storages || []}
        statusTip={
          <StatusTip
            isLoading={storageLoading}
            isError={storageHasError}
            error={storageError}
          />
        }
      />
    </>
  );
};

export default Overview;
