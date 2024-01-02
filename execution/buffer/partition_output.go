package buffer

import (
	"context"

	"google.golang.org/grpc"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/execution/model"
	"github.com/lindb/lindb/models"
	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	"github.com/lindb/lindb/spi"
)

type PartitionOutputBuffer struct {
	receivers []models.InternalNode
}

func NewPartitionOutputBuffer(receivers []models.InternalNode) OutputBuffer {
	return &PartitionOutputBuffer{
		receivers: receivers,
	}
}

// AddPage implements OutputBuffer
func (output *PartitionOutputBuffer) AddPage(page *spi.Page) {
	receiver := output.receivers[0]
	conn, err := grpc.Dial(receiver.Address(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := protoCommandV1.NewResultSetServiceClient(conn)
	_, err = client.ResultSet(context.TODO(), &protoCommandV1.ResultSetRequest{
		Payload: encoding.JSONMarshal(&model.TaskResultSet{
			TaskID: model.TaskID{RequestID: "1", ID: 1},
			NodeID: 6,
			Page:   page,
			NoMore: true,
		}),
	})

	if err != nil {
		panic(err)
	}
}
