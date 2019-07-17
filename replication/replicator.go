package replication

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/rpc/proto/common"
	"github.com/eleme/lindb/rpc/proto/storage"
)

const defaultBufferSize = 10

type Listener interface {
}

type Replicator interface {
	Replica(data *common.Replica)
}

type replicator struct {
	requestID int64
	writer    storage.WriteService_WriteClient
	buf       chan *common.Replica
}

func NewReplicator(ctx context.Context, conn *grpc.ClientConn, bufSize int) (Replicator, error) {
	size := bufSize
	if size <= 0 {
		size = defaultBufferSize
	}
	writer, err := storage.NewWriteServiceClient(conn).Write(ctx)
	if err != nil {
		return nil, err
	}

	r := &replicator{
		buf:    make(chan *common.Replica, size),
		writer: writer,
	}

	//todo
	go r.recv()
	go r.send()
	return r, nil
}

func (r *replicator) Replica(data *common.Replica) {
	r.buf <- data
}

func (r *replicator) recv() {
	resp, err := r.writer.Recv()
	//if err == io.EOF {
	//
	//}
	logger.GetLogger("replication").Info("receive:", logger.Error(err), logger.Any("resp", resp))
}

func (r *replicator) send() {
	for replica := range r.buf {
		reqID := r.requestID
		r.requestID++
		//	//err := r.writer.Send(&storage.WriteRequest{
		//	//	RequestID: reqID,
		//	//	Replica:   replica,
		//	//})
		//	//if err != nil {
		//	//	//TODO handle err
		//	//}
		logger.GetLogger("replication").Info("xxx", zap.Any("ss", reqID), zap.Any("dd", replica))
	}
}
