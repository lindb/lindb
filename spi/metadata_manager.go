package spi

import (
	"context"

	"github.com/lindb/common/pkg/encoding"
	"google.golang.org/grpc"

	"github.com/lindb/lindb/flow"
	protoMetaV1 "github.com/lindb/lindb/proto/gen/v1/meta"
)

type MetadataManager interface {
	GetTableMetadata(db, ns, name string) (*TableMetadata, error)
}

type metadataManager struct {
	nodeSelector flow.NodeSelector
}

func NewMetadataManager(nodeSelector flow.NodeSelector) MetadataManager {
	return &metadataManager{
		nodeSelector: nodeSelector,
	}
}

func (mgr *metadataManager) GetTableMetadata(db, ns, name string) (*TableMetadata, error) {
	partitions, err := mgr.nodeSelector.GetPartitions(db)
	if err != nil {
		return nil, err
	}
	schema := NewTableSchema()
	for node := range partitions {
		conn, err := grpc.Dial(node.Address(), grpc.WithInsecure())
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
