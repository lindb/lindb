package rpc

import (
	context "context"
	"fmt"

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

func (srv *MetaService) TableSchema(ctx context.Context, request *protoMetaV1.TableSchemaRequest) (*protoMetaV1.TableSchemaResponse, error) {
	fmt.Printf("mete request=%v\n", request)
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
		fmt.Printf("err1=%v\n", err)
		return nil, err
	}
	schema, err := database.MetaDB().GetSchema(metricID)
	if err != nil {
		fmt.Printf("err2=%v\n", err)
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
