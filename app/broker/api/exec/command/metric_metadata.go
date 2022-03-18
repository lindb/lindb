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

package command

import (
	"context"
	"strings"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// MetricMetadataCommand executes the metric metadata query.
func MetricMetadataCommand(ctx context.Context, deps *deps.HTTPDeps, param *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	metadataStmt := stmt.(*stmtpkg.MetricMetadata)
	if strings.TrimSpace(param.Database) == "" {
		return nil, constants.ErrDatabaseNameRequired
	}
	if metadataStmt.Limit == 0 || metadataStmt.Limit > constants.MaxSuggestions {
		// if limit =0 or > max suggestion items, need reset limit
		metadataStmt.Limit = constants.MaxSuggestions
	}
	//if metadataStmt.li
	return suggest(ctx, deps, param, metadataStmt)
}

// suggest executes metadata suggest query.
func suggest(ctx context.Context, deps *deps.HTTPDeps, param *models.ExecuteParam, request *stmtpkg.MetricMetadata) (interface{}, error) {
	metaDataQuery := deps.QueryFactory.NewMetadataQuery(ctx, param.Database, request)
	values, err := metaDataQuery.WaitResponse()
	if err != nil {
		return nil, err
	}
	switch request.Type {
	case stmtpkg.Field:
		// build field result model
		result := make(map[field.Name]field.Meta)
		fields := field.Metas{}
		for _, value := range values {
			err = encoding.JSONUnmarshal([]byte(value), &fields)
			if err != nil {
				return nil, err
			}
			for _, f := range fields {
				result[f.Name] = f
			}
		}
		// HistogramSum(sum), HistogramCount(sum), HistogramMin(min), HistogramMax(max) is visible
		// __bucket_{id}(HistogramField) is not visible for api,
		// underlying histogram data is only restricted access by user via quantile function
		// furthermore, we suggest some quantile functions for user in field names, such as quantile(0.99)
		var (
			resultFields []models.Field
			hasHistogram bool
		)
		for _, f := range result {
			if f.Type != field.HistogramField {
				resultFields = append(resultFields, models.Field{
					Name: string(f.Name),
					Type: f.Type.String(),
				})
			} else {
				hasHistogram = true
			}
		}
		//
		if hasHistogram {
			resultFields = append(resultFields,
				models.Field{Name: "quantile(0.99)", Type: field.HistogramField.String()},
				models.Field{Name: "quantile(0.95)", Type: field.HistogramField.String()},
				models.Field{Name: "quantile(0.90)", Type: field.HistogramField.String()},
			)
		}
		return &models.Metadata{
			Type:   request.Type.String(),
			Values: resultFields,
		}, nil
	default:
		return &models.Metadata{
			Type:   request.Type.String(),
			Values: values,
		}, nil
	}
}
