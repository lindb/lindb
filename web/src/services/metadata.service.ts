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
import { GET, POST } from "@src/utils";
import { Storage, Database } from "@src/models";
import { ApiPath } from "@src/constants";

/**
 * find all register storage cluster list.
 *  @return all storage cluster list
 */
export function findStorageList() {
  return GET<Storage[]>(ApiPath.StorageList);
}

export function createStorage(storage: Storage) {
  return POST<any>(ApiPath.Storage, storage);
}
/**
 * find all database config list.
 * @returns all database config list.
 */
export function findDatabaseList() {
  return GET<Database[]>(ApiPath.DatabaseList);
}

export function metaExplore() {
  return GET<any>(ApiPath.MetaExplore);
}

export function exploreRepoData(params: any) {
  return GET<any>(ApiPath.RepoExplore, params);
}
