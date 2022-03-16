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
export type Master = {
  node?: Node;
  electTime?: number;
};

export type Node = {
  id?: number;
  hostIp?: string;
  hostName?: string;
  grpcPort?: number;
  httpPort?: number;
  version?: string;
  onlineTime?: string;
};

export enum ShardStateType {
  UnknownShard = 0,
  NewShard = 1,
  OnlineShard = 2,
  OfflineShard = 3,
  NonExistentShard = 4,
}

export type StorageState = {
  name?: string;
  liveNodes?: { [nodeId: string]: Node }; // node id=>node
  shardStates?: {
    [database: string]: {
      [shardId: string]: {
        id?: number;
        leader?: number;
        replica: { [propName: string]: number[] }; // replicas=> node ids
        state?: ShardStateType;
      };
    };
  }; // database's nmae => shard state
};

export type StateMetric = {
  [name: string]: {
    fields: { name?: string; type?: string; value?: number }[];
    tags: { [key: string]: string };
  }[];
};

export type ReplicaState = {
  [node: string]: {
    shardId: number;
    familyTime: string;
    leader: number;
    append: number; //next write index
    replicators: {
      replicator: string; //node id
      consume: number; //next consume idx
      ack: number;
      pending: number;
    }[];
  }[];
};
