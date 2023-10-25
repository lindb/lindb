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
  DatabaseView,
  NodeView,
  StatusTip,
  StorageView,
} from "@src/components";
import { SQL, StateMetricName } from "@src/constants";
import { UIContext } from "@src/context/UIContextProvider";
import { useParams } from "@src/hooks";
import { ExecService } from "@src/services";
import { StateKit } from "@src/utils";
import { useQuery } from "@tanstack/react-query";
import * as _ from "lodash-es";
import React, { useContext, useEffect } from "react";

const StorageOverview: React.FC = () => {
  const { name } = useParams(["name"]);
  const {
    isLoading,
    isInitialLoading,
    isFetching,
    isError,
    error,
    data,
    refetch,
  } = useQuery(
    ["show_alive_storage", name],
    async () => {
      return ExecService.exec<any[]>({ sql: SQL.ShowStorageAliveNodes });
    },
    { enabled: !_.isEmpty(name) }
  );
  const storages = StateKit.getStorageState(data, name || "");
  const { locale } = useContext(UIContext);
  const { StorageView: StorageViewCmp } = locale;

  // reload when change language
  useEffect(() => {
    refetch();
  }, [locale, refetch]);

  if (isError || isLoading || isInitialLoading || isFetching) {
    return (
      <StatusTip
        style={{ marginTop: 150 }}
        isLoading={isLoading || isInitialLoading || isFetching}
        isError={isError}
        error={error}
      />
    );
  }

  return (
    <>
      <StorageView name={name as string} storages={storages || []} />
      <NodeView
        showNodeId
        title={StorageViewCmp.liveNodes}
        nodes={_.orderBy(
          _.values(_.get(storages, "[0].liveNodes", {})),
          ["id"],
          ["asc"]
        )}
        sql={`show storage metric where storage='${name}' and metric in ('${StateMetricName.CPU}','${StateMetricName.Memory}')`}
        style={{ marginTop: 12, marginBottom: 12 }}
      />
      <DatabaseView
        title={StorageViewCmp.databaseList}
        liveNodes={_.get(storages, "[0].liveNodes", {})}
        storage={_.get(storages, "[0]", {})}
      />
    </>
  );
};

export default StorageOverview;
