package spi

import (
	"context"

	"github.com/lindb/common/pkg/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lindb/lindb/flow"
	protoMetaV1 "github.com/lindb/lindb/proto/gen/v1/meta"
)

type CreateTableHandle func(db, ns, name string) TableHandle

var createTableHandleFn = make(map[TableKind]CreateTableHandle)

func RegisterCreateTableHandleFn(kind TableKind, fn CreateTableHandle) {
	createTableHandleFn[kind] = fn
}

type MetadataManager interface {
	GetTableMetadata(db, ns, name string) (*TableMetadata, error)
	GetTableHandle(db, ns, name string) TableHandle
}

type metadataManager struct {
	nodeSelector flow.NodeSelector
}

func NewMetadataManager(nodeSelector flow.NodeSelector) MetadataManager {
	return &metadataManager{
		nodeSelector: nodeSelector,
	}
}

func (mgr *metadataManager) GetTableHandle(db, ns, name string) TableHandle {
	// FIXME: set table kind
	fn, ok := createTableHandleFn[MetricTable]
	if !ok {
		panic("create table handle fn not exist")
	}
	return fn(db, ns, name)
}

func (mgr *metadataManager) GetTableMetadata(db, ns, name string) (*TableMetadata, error) {
	partitions, err := mgr.nodeSelector.GetPartitions(db)
	if err != nil {
		return nil, err
	}
	schema := NewTableSchema()
	for node := range partitions {
		conn, err := grpc.Dial(node.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := protoMetaV1.NewMetaServiceClient(conn)
		resp, err := client.TableSchema(context.TODO(), &protoMetaV1.TableSchemaRequest{
			Database:  db,
			Namespace: ns,
			Table:     name,
		})
		if err != nil {
			return nil, err
		}
		tableSchema := &TableSchema{}
		if err = encoding.JSONUnmarshal(resp.Payload, tableSchema); err != nil {
			return nil, err
		}
		schema.AddColumns(tableSchema.Columns)
	}
	return &TableMetadata{
		Schema:     schema,
		Partitions: partitions,
	}, nil
}
