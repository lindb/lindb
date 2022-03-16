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
import { useAliveState } from "@src/hooks";
import { StorageState } from "@src/models";
import * as _ from "lodash-es";
import { useEffect, useState } from "react";

const aliveStorage = SQL.ShowStorageAliveNodes;

export function useStorage(name?: string) {
  const [storages, setStorages] = useState<StorageState[]>();
  const { loading, aliveState } = useAliveState(aliveStorage);

  useEffect(() => {
    const fetchStorageState = () => {
      try {
        let storages = aliveState;
        if (name) {
          console.log("storages", storages, name);
          const idx = _.findIndex(storages, { name: name });
          console.log("storages", storages, name, idx);
          storages = idx >= 0 ? _.pullAt(storages, [idx]) : [];
        }
        (storages || []).map((storage: StorageState) => {
          const liveNodes = _.get(storage, "liveNodes", {});
          const databases = _.get(storage, "shardStates", {});
          const stats = {
            numOfDatabase: 0,
            totalReplica: 0,
            availableReplica: 0,
            unavailableReplica: 0,
            liveNodes: _.keys(liveNodes).length,
            deadNodes: [] as number[],
          };
          _.set(storage, "stats", stats);
          _.mapValues(databases, function (db: any) {
            stats.numOfDatabase++;
            _.mapValues(db, function (shard: any) {
              const replicas = _.get(shard, "replica.replicas", []);
              stats.totalReplica += replicas.length;
              replicas.map((nodeId: number) => {
                if (_.has(liveNodes, nodeId)) {
                  stats.availableReplica++;
                } else {
                  stats.unavailableReplica++;
                  stats.deadNodes.push(nodeId);
                }
              });
              console.log("shard....", shard, replicas);
            });
          });
          stats.deadNodes = _.uniq(stats.deadNodes);
        });
        console.log("storage...", storages);
        setStorages(storages);
      } catch (err) {
        console.log(err);
      }
    };
    fetchStorageState();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [name, aliveState]);

  return { loading, storages };
}
