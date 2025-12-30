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

package streaming

import (
	"fmt"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/streaming/transfer"
)

type DataSource interface {
	Initialize()
	Name() string
	Produce(data []byte)
}

type dataSource struct {
	db models.Database

	transfer transfer.Transfer
}

func NewDataSource(db models.Database) DataSource {
	return &dataSource{
		db: db,
	}
}

func (d *dataSource) Initialize() {
	d.transfer = transfer.GetTransfer(d.db.Option.Engine)
}

func (d *dataSource) Name() string {
	return d.db.Name
}

func (d *dataSource) Produce(data []byte) {
	page, err := d.transfer.ToPage(data)
	if err != nil {
		fmt.Println("produce trace data error:", err)
	}
	fmt.Println("produce trace data:", string(encoding.JSONMarshal(page)))
}
