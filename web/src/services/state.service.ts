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
import { ApiPath } from "@src/constants";
import { ReplicaState, StateMetric } from "@src/models";
import { GET } from "@src/utils";

export function exploreState(params: any) {
  return GET<StateMetric>(ApiPath.StateExplore, params);
}

export function getReplicaState(params: any) {
  return GET<ReplicaState>(ApiPath.ReplicaState, params);
}
