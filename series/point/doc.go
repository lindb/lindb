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

/*
Note: This package is heavily referenced to the InfluxDB project,
      as the design goal is to provide an InfluxDB-compatible line-protocol for LinDB.

InfluxDB line protocol syntax
Reference: https://docs.influxdata.com/influxdb/v1.7/write_protocols/line_protocol_tutorial/

weather,location=us-midwest temperature=82 1465839830100400200
  |    -------------------- --------------  |
  |             |             |             |
  |             |             |             |
+-----------+--------+-+---------+-+---------+
|measurement|,tag_set| |field_set| |timestamp|
+-----------+--------+-+---------+-+---------+



In LinDB, fields are different from InfluxDB,
there are multi types of fields, such as SumField, HistogramField, etc.

Mapping Example for Metric Type
Counter -> SumField
Gauge -> SumField(total) + SumField(count)


The field type is appended to the field-name,
It is a 4 byte slice startswith an underline and the upper case abbreviation of field-type,
such as `_SUM`(sum field), `_MIN`(min field), `_MAX`(max field), `_SMY`(SummaryField)
The enhanced example is shown below

weather,location=us-midwest temperature_SUM=82 1465839830100400200
  |    -------------------- --------------     |
  |             |             |               |
  |             |             |              |
+-----------+--------+-+---------+-+---------+
|metric-name|,tag_set| |field_set| |timestamp|
+-----------+--------+-+---------+-+---------+


*/
package point
