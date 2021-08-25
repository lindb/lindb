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

package replica

/**
Channel represents a persistent, replicable storage for a data shard in Broker.
One data shard may be copied to many storage for fault-tolerance.
Replicator handles the details of copping data shard to a target Storage node.

                                                             Storage
                              +--------------+               +-------------------+
                              |Replicator    |               |  StorageService   |
                              |              |  Write        |                   |
Broker                        |Target        +-------------->+  SequenceManager  |
+----------------------+      |Cluster       |  Stream RPC   |  HeadSeq          |
| Channel              |      |Database      <---------------+  AckSeq           |
|                      +------+ShardID       |  Ack Reset    |                   |
| Cluster              |      |Pending       |               |                   |
| Database             |      |Stop          |               |                   |
| ShardID              |      |              |               |                   |
| Write                |      +--------------+               +-------------------+
| GetOrCreateReplicator|
| Targets              |
|                      |
|                      |
+----------------------+      Replicators......

*/
