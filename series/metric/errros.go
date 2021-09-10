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

package metric

import (
	"errors"
	"fmt"
)

var (
	ErrTooManyTags = errors.New("too_many_tags")
	// ErrBadMetricPBFormat represents write bad pb format
	ErrBadMetricPBFormat = errors.New("bad metric proto")
	ErrMetricPBNilMetric = fmt.Errorf("%w, metric is nil", ErrBadMetricPBFormat)
	// ErrMetricPBEmptyMetricName represents metric name is empty when write data
	ErrMetricPBEmptyMetricName = fmt.Errorf("%w, metric name is empty", ErrBadMetricPBFormat)
	ErrMetricEmptyTagKeyValue  = fmt.Errorf("%w tag key value is empty", ErrBadMetricPBFormat)
	// ErrMetricPBEmptyField represents field is empty when write data
	ErrMetricPBEmptyField = fmt.Errorf("%w, fields are empty", ErrBadMetricPBFormat)
	// ErrMetricEmptyFieldName represents that field-name is empty in pb structure
	ErrMetricEmptyFieldName = fmt.Errorf("%w, field name is empty", ErrBadMetricPBFormat)
	// ErrMetricNanField represents field value is not a number
	ErrMetricNanField = fmt.Errorf("%w, field is not a number", ErrBadMetricPBFormat)
	// ErrMetricInfField represents field value is infinity, positive or negative
	ErrMetricInfField = fmt.Errorf("%w, field is infinity", ErrBadMetricPBFormat)
)
