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

package state

import (
	"fmt"

	etcdcliv3 "go.etcd.io/etcd/client/v3"
)

var (
	// ErrWatchFailed indicates the watch failed.
	ErrWatchFailed = fmt.Errorf("etcd watch returns a nil chan")
	// ErrNoKey indicates the key does not exist.
	ErrNoKey = fmt.Errorf("etcd has no such key")
	// ErrTxnFailed indicates the txn failed.
	ErrTxnFailed = fmt.Errorf("role changed or target revision mismatch")
	// ErrTxnConvert transaction covert failed.
	ErrTxnConvert = fmt.Errorf("cannot covert etcd transaction")
)

// TxnErr converts txn response and error into one error.
func TxnErr(resp *etcdcliv3.TxnResponse, err error) error {
	if err != nil {
		return err
	}
	if !resp.Succeeded {
		return ErrTxnFailed
	}
	return nil
}
