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
import { NodeView, StorageView, DatabaseView } from "@src/components";
import { StateMetricName, StateRoleName } from "@src/constants";
import { useStorage } from "@src/hooks";
import * as _ from "lodash-es";
import { URLStore } from "@src/stores";
import React from "react";

export default function StorageOverview() {
  const name = URLStore.params.get("name");
  const { loading, storages } = useStorage(name as string);
  return (
    <>
      {!loading && (
        <>
          <StorageView
            name={name as string}
            storages={storages || []}
            loading={loading}
          />
          <NodeView
            title="Live Nodes"
            loading={loading}
            nodes={_.values(_.get(storages, "[0].liveNodes", {}))}
            stateParams={{
              names: [StateMetricName.CPU, StateMetricName.Memory],
              role: StateRoleName.Storage,
              storageName: name,
            }}
            style={{ marginTop: 12, marginBottom: 12 }}
          />
          <DatabaseView
            liveNodes={_.get(storages, "[0].liveNodes", {})}
            storage={_.get(storages, "[0]", {})}
            loading={loading}
          />
        </>
      )}
    </>
  );
}
