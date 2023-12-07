// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package memdb

// ┌────────────────────────────────┐
// │Metric ID -> Metric Index(Map)  │
// └──────────────┬─────────────────┘     ┌───────────────────────────────────┐
//                │                       │           Metric Store            │
//   Metric Index │                       │ ┌───────────────────────────────┐ │
//                │                       │ │Tags Hash -> Mem Series ID(Map)│ │
//                ▼                       │ └───────────────────────────────┘ │
//       ┌───────────────────┐            │                                   │   Map Ref(Mem Series ID->Field Store)
//       │Metric Store(Array)│ ────────── │ ┌───────────────────────────────┐ │        ┌──────────────────────┐
//       └───────────────────┘            │ │Series ID -> Mem Series ID(Map)│ ├────┬──►│Field Store Map(Array)│
//                ▲                       │ └───────────────────────────────┘ │    │   └──────────────────────┘
//                │                       │                                   │    │
//   Metric Index │                       │ ┌───────────────────────────────┐ │ Array Ref(Mem Field Index)
//                │                       │ │Fields(Field Mem Index/ID)     │ │
// ┌──────────────┴─────────────────┐     │ └───────────────────────────────┘ │
// │Metric Hash -> Metric Index(Map)│     │                                   │
// └────────────────────────────────┘     └───────────────────────────────────┘
