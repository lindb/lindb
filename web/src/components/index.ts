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

// layout component
export { default as Footer } from "@src/components/layout/Footer";
export { default as Header } from "@src/components/layout/Header";
export { default as SiderMenu } from "@src/components/layout/SiderMenu";

// commom component
export { default as TimePicker } from "@src/components/common/TimePicker";
export { default as StatusTip } from "@src/components/common/StatusTip";
export { default as SimpleStatusTip } from "@src/components/common/SimpleStatusTip";
export { default as LazyLoad } from "@src/components/common/LazyLoad";
export { default as Icon } from "@src/components/common/Icon";

// form component
export { default as LinSelect } from "@src/components/form/Select";

// chart component
export * from "@src/components/chart/Chart";

// data component
export { default as MetadataSelect } from "@src/components/data/MetadataSelect";
export { default as TagFilterSelect } from "@src/components/data/TagFilterSelect";
export { default as TagValueSelect } from "@src/components/data/TagValueSelect";
export { default as ExplainStatsView } from "@src/components/data/ExplainStatsView";

// state component
export { default as MasterView } from "@src/components/state/MasterView";
export { default as NodeView } from "@src/components/state/NodeView";
export { default as StorageView } from "@src/components/state/StorageView";
export { default as DiskUsageView } from "@src/components/state/DiskUsageView";
export { default as CapacityView } from "@src/components/state/CapacityView";
export { default as DatabaseView } from "@src/components/state/DatabaseView";
export { default as ReplicaView } from "@src/components/state/ReplicaView";
export { default as StorageStatusView } from "@src/components/state/StorageStatusView";
export { default as MemoryDatabaseView } from "@src/components/state/MemoryDatabaseView";
