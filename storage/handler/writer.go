package handler

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	streamIO "github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/rpc/proto/storage"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/tsdb"
)

// Writer implements the stream write service.
type Writer struct {
	storageService service.StorageService
	sm             replication.SequenceManager
	logger         *logger.Logger
}

// NewWriter returns a new Writer.
func NewWriter(storageService service.StorageService, sm replication.SequenceManager) *Writer {
	return &Writer{
		storageService: storageService,
		sm:             sm,
		logger:         logger.GetLogger("storage", "Writer"),
	}
}

func (w *Writer) Reset(ctx context.Context, req *storage.ResetSeqRequest) (*storage.ResetSeqResponse, error) {
	logicNode, err := getLogicNodeFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sequence, err := w.getSequence(req.Database, req.ShardID, *logicNode)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	sequence.SetHeadSeq(req.Seq)

	return &storage.ResetSeqResponse{}, nil
}

func (w *Writer) Next(ctx context.Context, req *storage.NextSeqRequest) (*storage.NextSeqResponse, error) {
	logicNode, err := getLogicNodeFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sequence, err := w.getSequence(req.Database, req.ShardID, *logicNode)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &storage.NextSeqResponse{Seq: sequence.GetHeadSeq()}, nil
}

// Write handles the stream write request.
func (w *Writer) Write(stream storage.WriteService_WriteServer) error {
	database, shardID, logicNode, err := parseCtx(stream.Context())
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	sequence, err := w.getSequence(database, shardID, *logicNode)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	shard := w.storageService.GetShard(database, shardID)
	if shard == nil {
		return status.Errorf(codes.NotFound, "shard %d for database %s not exists", shardID, database)
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			w.logger.Error("error", logger.Error(err))
			return status.Error(codes.Internal, err.Error())
		}

		if len(req.Replicas) == 0 {
			continue
		}

		// nextSeq means the sequence replica wanted
		for _, replica := range req.Replicas {
			seq := replica.Seq

			hs := sequence.GetHeadSeq()
			if hs != seq {
				// reset to headSeq
				return status.Errorf(codes.OutOfRange, "seq num not match replica:%d, storage:%d", seq, hs)
			}

			w.handleReplica(shard, replica)

			sequence.SetHeadSeq(hs + 1)

		}

		resp := &storage.WriteResponse{
			CurSeq: sequence.GetHeadSeq() - 1,
		}

		// add acked seq if synced
		if sequence.Synced() {
			resp.Ack = &storage.WriteResponse_AckSeq{AckSeq: sequence.GetAckSeq()}
		}

		if err := stream.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

func (w *Writer) handleReplica(shard tsdb.Shard, replica *storage.Replica) {
	reader := streamIO.NewReader(replica.Data)
	for !reader.Empty() {
		bytesLen := reader.ReadUvarint32()

		bytes := reader.ReadBytes(int(bytesLen))

		if err := reader.Error(); err != nil {
			w.logger.Error("read metricList bytes from replica", logger.Error(err))
			break
		}

		var metricList field.MetricList
		err := metricList.Unmarshal(bytes)
		if err != nil {
			w.logger.Error("unmarshal metricList", logger.Error(err))
			continue
		}

		//todo DEBUG level
		w.logger.Info("receive metricList", logger.Any("metricList", metricList))

		//TODO write metric, need handle panic
		for _, metric := range metricList.Metrics {
			err = shard.Write(metric)
		}
		if err != nil {
			w.logger.Error("write metric", logger.Error(err))
			continue
		}
	}
}

func getLogicNodeFromCtx(ctx context.Context) (*models.Node, error) {
	return rpc.GetLogicNodeFromContext(ctx)
}

func parseCtx(ctx context.Context) (database string, shardID int32, logicNode *models.Node, err error) {
	logicNode, err = rpc.GetLogicNodeFromContext(ctx)
	if err != nil {
		return
	}

	database, err = rpc.GetDatabaseFromContext(ctx)
	if err != nil {
		return
	}

	shardID, err = rpc.GetShardIDFromContext(ctx)
	return
}

func (w *Writer) getSequence(database string, shardID int32, logicNode models.Node) (replication.Sequence, error) {
	sequence, ok := w.sm.GetSequence(database, shardID, logicNode)
	if !ok {
		var err error
		sequence, err = w.sm.CreateSequence(database, shardID, logicNode)
		if err != nil {
			return nil, err
		}
	}
	return sequence, nil
}
