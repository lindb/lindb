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

// Package operator provides the core operators for the execution pipeline.
package operator

// A distributed SQL query engine designed for time series(ref:https://trino.io/).
// It supports a variety of SQL operations, similar to traditional SQL databases,
// but also includes optimizations and features suited for time series environments.
// Here are some key operators and concepts in LinDB:
//
// 1. TableScan: reads data from a data source(metric/log etc.)
// 2. Projection: applies transformations to the input data, suchs computing expressions based on input columns
// 3. Filter: filters rows that do not meet the condition specified
// 4. Aggregation: performs calculations across a set of rows that are grouped together based on one or more columns(sum/avg/count etc.)
// 5. Join: combines rows from two or more tables based on a related column between them
// 6: Sort: order the data based on one or more columns(asc/desc)
// 7: Limit: limits the number of rows returned by a query
// 8: Window Function?
// 9: TopN?
// 10: Hash?
