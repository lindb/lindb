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

package metadb

/*
Metric metadata storage format using bbolt:

+-------------+
| root Bucket |
+------+------+
       |       +----------------+
       +------>+ ns root bucket +-----> metric id global sequence in bucket
       |       +-------+--------+
       |               |         +-------------+
       |               +-------->+ ns bucket 1 |
       |               |         +------+------+
       |               |                |        +-------------------+
       |               |                +------->+ metric name 1->id |
       |               |                |        +-------------------+
       |               |                |        +-------------------+
       |               |                +------->+ metric name n->id |
       |               |                         +-------------------+
       |               |         +-------------+
       |               +-------->+ ns bucket 2 |
       |               |         +------+------+
       |               |                |        +-------------------+
       |               |                +------->+ metric name 1->id |
       |               |                |        +-------------------+
       |               |                |        +-------------------+
       |               |                +------->+ metric name n->id |
       |               |                         +-------------------+
       |               |         +-------------+
       |               +-------->+ ns bucket n |
       |                         +-------------+
       |      +--------------------+
       +----->+ metric root bucket +------> tag id global sequence in bucket
              +---------+----------+
                        |            +-----------------+
                        +----------->+ metric bucket 1 |
                        |            +--------+--------+
                        |                     |         +----------------+
                        |                     +-------->+ tag key bucket |
                        |                     |         +-------+--------+
                        |                     |                 |          +---------------+
                        |                     |                 +--------->+ tag key 1->id |
                        |                     |                 |          +---------------+
                        |                     |                 |          +---------------+
                        |                     |                 +--------->+ tag key n->id |
                        |                     |                            +---------------+
                        |                     |         +--------------+
                        |                     +-------->+ field bucket +-------> field id sequence of metric level
                        |                               +------+-------+
                        |                                      |        +------------------------+
                        |                                      +------->+ field name->field meta |
                        |                                      |        +------------------------+
                        |                                      |        +------------------------+
                        |                                      +------->+ field name->field meta |
                        |                                               +------------------------+
                        |            +-----------------+
                        +----------->+ metric bucket n |
                                     +--------+--------+
*/
