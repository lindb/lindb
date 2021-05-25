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

package parallel

import "errors"

var errUnmarshalPlan = errors.New("unmarshal physical plan error")
var errUnmarshalQuery = errors.New("unmarshal query statement error")
var errUnmarshalSuggest = errors.New("unmarshal metadata suggest statement error")
var errWrongRequest = errors.New("not found task of current node from physical plan")
var errNoSendStream = errors.New("not found send stream")
var errTaskSend = errors.New("send task request error")
var errNoDatabase = errors.New("not found database")
