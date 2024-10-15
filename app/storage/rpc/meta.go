package rpc

import (
	context "context"

	commonConstants "github.com/lindb/common/constants"
	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
	protoMetaV1 "github.com/lindb/lindb/proto/gen/v1/meta"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/tsdb"
)

type MetaService struct {
	engine tsdb.Engine
	logger logger.Logger
}

func NewMetaService(engine tsdb.Engine) protoMetaV1.MetaServiceServer {
	return &MetaService{
		engine: engine,
		logger: logger.GetLogger("RPC", "resultSet"),
	}
}

func (srv *MetaService) SuggestNamespace(ctx context.Context, request *protoMetaV1.SuggestRequest) (*protoMetaV1.SuggestResponse, error) {
	database, ok := srv.engine.GetDatabase(request.Database)
	if !ok {
		return nil, constants.ErrDatabaseNotFound
	}
	namespaces, err := database.MetaDB().SuggestNamespace(request.Namespace, int(request.Limit))
	if err != nil {
		return nil, err
	}
	return &protoMetaV1.SuggestResponse{Values: namespaces}, nil
}

func (srv *MetaService) SuggestTable(ctx context.Context, request *protoMetaV1.SuggestRequest) (*protoMetaV1.SuggestResponse, error) {
	database, ok := srv.engine.GetDatabase(request.Database)
	if !ok {
		return nil, constants.ErrDatabaseNotFound
	}
	if !ok {
		return nil, constants.ErrDatabaseNotFound
	}
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

func (srv *MetaService) TableSchema(ctx context.Context, request *protoMetaV1.TableSchemaRequest) (*protoMetaV1.TableSchemaResponse, error) {
	database, ok := srv.engine.GetDatabase(request.Database)
	if !ok {
		return nil, constants.ErrDatabaseNotFound
	}
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
	tableSchema := types.NewTableSchema()
	for _, tagKey := range schema.TagKeys {
		tableSchema.AddColumn(types.ColumnMetadata{Name: tagKey.Key, DataType: types.DTString})
	}
	for _, field := range schema.Fields {
		tableSchema.AddColumn(types.ColumnMetadata{
			Name:     field.Name.String(),
			DataType: types.DTTimeSeries,
			AggType:  field.Type.AggregateType(),
		})
	}
	return &protoMetaV1.TableSchemaResponse{
		Payload: encoding.JSONMarshal(tableSchema),
	}, nil
}
