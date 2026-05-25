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

package rpc

import (
	context "context"

	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	commonConstants "github.com/lindb/common/constants"
	"github.com/lindb/common/field"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
	protoMetaV1 "github.com/lindb/lindb/proto/gen/v1/meta"
	"github.com/lindb/lindb/storage"
	"github.com/lindb/lindb/storage/metric"
)

type MetaService struct {
	protoMetaV1.UnimplementedMetaServiceServer

	engine storage.Engine
	logger logger.Logger
}

func NewMetaService(engine storage.Engine) protoMetaV1.MetaServiceServer {
	return &MetaService{
		engine: engine,
		logger: logger.GetLogger("RPC", "resultSet"),
	}
}

func (srv *MetaService) SuggestNamespace(ctx context.Context,
	request *protoMetaV1.SuggestRequest,
) (*protoMetaV1.SuggestResponse, error) {
	db, ok := srv.engine.GetDatabase(request.Database)
	if !ok {
		return nil, constants.ErrDatabaseNotFound
	}
	database := db.(*metric.Database)
	namespaces, err := database.MetaDB().SuggestNamespace(request.Namespace, int(request.Limit))
	if err != nil {
		return nil, err
	}
	return &protoMetaV1.SuggestResponse{Values: namespaces}, nil
}

func (srv *MetaService) SuggestTable(ctx context.Context,
	request *protoMetaV1.SuggestRequest,
) (*protoMetaV1.SuggestResponse, error) {
	db, ok := srv.engine.GetDatabase(request.Database)
	if !ok {
		return nil, constants.ErrDatabaseNotFound
	}
	database := db.(*metric.Database)
	namespace := commonConstants.DefaultNamespace
	if request.Namespace != "" {
		namespace = request.Namespace
	}
	metrics, err := database.MetaDB().SuggestMetrics(namespace, request.Table, int(request.Limit))
	if err != nil {
		return nil, err
	}
	return &protoMetaV1.SuggestResponse{Values: metrics}, nil
}

func (srv *MetaService) TableSchema(ctx context.Context,
	request *protoMetaV1.TableSchemaRequest,
) (*protoMetaV1.TableSchemaResponse, error) {
	db, ok := srv.engine.GetDatabase(request.Database)
	if !ok {
		return nil, constants.ErrDatabaseNotFound
	}
	database := db.(*metric.Database)
	namespace := commonConstants.DefaultNamespace
	if request.Namespace != "" {
		namespace = request.Namespace
	}
	metricID, err := database.MetaDB().GetMetricID(namespace, request.Table)
	if err != nil {
		return nil, err
	}
	schema, err := database.MetaDB().GetSchema(metricID)
	if err != nil {
		return nil, err
	}
	// TODO: return schema for metric engine
	var fields []arrow.Field
	for _, tagKey := range schema.TagKeys {
		fields = append(fields, arrow.Field{Name: tagKey.Key, Type: arrow.BinaryTypes.String})
	}
	getDataType := func(fieldType field.Type) arrow.DataType {
		switch fieldType {
		case field.Exemplar:
			return larrow.ExtensionTypes.Exemplar
		case field.Sum:
			return larray.NewAggregationType(larray.Sum)
		case field.Min:
			return larray.NewAggregationType(larray.Min)
		case field.Max:
			return larray.NewAggregationType(larray.Max)
		case field.Last:
			return larray.NewAggregationType(larray.Last)
		case field.First:
			return larray.NewAggregationType(larray.First)
		default:
			return larray.NewAggregationType(larray.Sum)
		}
	}
	// Emit all physical fields as-is; histogram bucket/stat fields are plain sum/min/max.
	// The naming convention (<histoName>.__bucket_*, <histoName>.__sum etc.) carries all
	// semantic information needed by the query and aggregation layers.
	for _, f := range schema.Fields {
		fields = append(fields, arrow.Field{Name: f.Name.String(), Type: getDataType(f.Type)})
	}

	// Timestamp is an implicit dimension, not a selectable data column.
	fields = append(fields, arrow.Field{
		Name:     constants.TimestampColumnName,
		Type:     arrow.FixedWidthTypes.Timestamp_ns,
		Metadata: arrow.MetadataFrom(map[string]string{"hidden": "true"}),
	})
	data, err := larrow.MarshalSchema(arrow.NewSchema(fields, nil))
	if err != nil {
		return nil, err
	}
	return &protoMetaV1.TableSchemaResponse{
		Payload: data,
	}, nil
}
