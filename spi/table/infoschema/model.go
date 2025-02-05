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

package infoschema

import (
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi"
)

func init() {
	// register table handle
	encoding.RegisterNodeType(TableHandle{})
	spi.RegisterCreateTableFn(spi.InfoSchema, func(db, ns, name string) spi.TableHandle {
		return &TableHandle{
			Table: name,
		}
	})
}

type TableHandle struct {
	Table string `json:"table"`
}

func (t *TableHandle) SetTimeRange(timeRange timeutil.TimeRange) {}

func (t *TableHandle) GetTimeRange() timeutil.TimeRange {
	return timeutil.TimeRange{}
}

func (t *TableHandle) SetInterval(interval timeutil.Interval) {}

func (t *TableHandle) GetInterval() timeutil.Interval {
	return timeutil.Interval(0)
}

// Kind returns the datasource kind.
func (t *TableHandle) Kind() spi.DatasourceKind {
	return spi.InfoSchema
}

// String returns the table info of information schema.
func (t *TableHandle) String() string {
	return fmt.Sprintf("%s.%s", constants.InformationSchema, t.Table)
}

type Snippet struct {
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
	Type     string `yaml:"type"`
}

type Function struct {
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
	Type     string `yaml:"type"`
}
