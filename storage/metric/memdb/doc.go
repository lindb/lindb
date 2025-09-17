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

// Memory Structure
//
// +-----------------------------------------------------+
// |                                                     |
// | Metric Meta Database(Database Level)                |           +--------------------------+
// |                                                     |           |                          |
// |                                                     |           | Metric Store             |
// | Metric Hash(namespace+metric name) =>  Metric Store-+---------->|                          |
// |                                                     |           |                          |
// | Metric ID => Metric Hash                            |           | Field Name => Field Meta |
// +-----------------------------------------------------+           +--------------------------+
//
// +--------------------------------------+
// |                                      |
// |  Metric Index Database(Shard Level)  |            +------------------------------------------+
// |                                      |            |                                          |
// |                                      |            | Time Series Index                        |
// |  Time Series Sequence(Memory Level)  |            |                                          |
// |                                      |            |                                          |
// |  Metric Hash => Time Series Index----+----------->| Tags Hash => Memory Time Series ID       |
// |                                      |            |                                          |
// +--------------------------------------+            | Time Series ID => Memory Time Series ID  |
//                                                     +------------------------------------------+
//
// +--------------------------------------------------------+
// |                                                        |
// |            Memory Database(Family Level)               |
// |                                                        |
// |                                                        |
// |                Memory Time Series ID                   |
// |                                                        |
// |                 |                 |                    |
// |                 |                 |                    |
// | Write Buffer <--+                 +--> Compress Buffer |
// |                                                        |
// +--------------------------------------------------------+
//
//
// Write Flow
//                                              3.write metric data
//                          +--------------------+           +-----------------------+
//                          | Data Family Writer +---------->|    Memory Database    |
//                          |    (goroutine)     |           |(Write/Compress buffer)|
//                          +--------------------+           +-----------------------+
//                            ^               ^
//   1.lookup mem metric meta |               | 2.lookup mem time series
//                            v               v
//            +---------------------+    +----------------------+
//            |Memroy Metric Meta DB|    |Memory Metric Index DB|
//            +---------------------+    +----------------------+
//                      ^                            ^ 1.watch event
//        1.watch event | 3.index metric meta        |
//                      v                            v 3. index time series
//              +------------------+    +------------------------+
//              |Metric Meta Lookup|    |Metric Tag/Sereis Lookup|
//              |   (goroutine)    |    |     (goroutine)        |
//              +------------------+    +------------------------+
//                    ^                               ^
//                    |                               |
// 2.lookup metric meta|                               | 2.lookup time series
//                    |     +--------------------+    |
//                    |     |Metric Meta/Index DB|    |
//                    +---->|     (Persist)      |<---+
//                          +--------------------+
