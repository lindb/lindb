package handler

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/rpc/proto/storage"
	"github.com/lindb/lindb/service"
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
		return nil, err
	}

	sequence, err := w.getSequence(req.Database, req.ShardID, *logicNode)
	if err != nil {
		return nil, err
	}

	sequence.SetHeadSeq(req.Seq)

	return &storage.ResetSeqResponse{}, nil
}

func (w *Writer) Next(ctx context.Context, req *storage.NextSeqRequest) (*storage.NextSeqResponse, error) {
	logicNode, err := getLogicNodeFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	sequence, err := w.getSequence(req.Database, req.ShardID, *logicNode)
	if err != nil {
		return nil, err
	}

	return &storage.NextSeqResponse{Seq: sequence.GetHeadSeq()}, nil
}

// Write handles the stream write request.
func (w *Writer) Write(stream storage.WriteService_WriteServer) error {
	database, shardID, logicNode, err := parseCtx(stream.Context())
	if err != nil {
		return err
	}

	sequence, ok := w.sm.GetSequence(database, shardID, *logicNode)
	if !ok {
		sequence, err = w.sm.CreateSequence(database, shardID, *logicNode)
		if err != nil {
			return err
		}
	}

	shard := w.storageService.GetShard(database, shardID)
	if shard == nil {
		return fmt.Errorf("shard %d for database %s not exists", shardID, database)
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			w.logger.Error("error", logger.Error(err))
			return err
		}

		if len(req.Replicas) == 0 {
			continue
		}

		// nextSeq means the sequence of replica wanted
		for _, rep := range req.Replicas {
			seq := rep.Seq

			hs := sequence.GetHeadSeq()
			if hs != seq {
				// reset to headSeq
				return errors.New("seq num not match")
			}

			sequence.SetHeadSeq(hs + 1)

			metric := &field.Metric{}
			//TODO need modify
			err := metric.Unmarshal(rep.Data)
			if err != nil {
				w.logger.Error("unmarshal metric", logger.Error(err))
				continue
			}
			w.logger.Info("receive metric", logger.Any("metric", rep.Data))

			//TODO write metric, need handle panic
			//err = shard.Write(metric)
			//if err != nil {
			//	logger.GetLogger("write").Error("write metric", logger.Error(err))
			//	continue
			//}
		}

		resp := &storage.WriteResponse{
			CurSeq: sequence.GetHeadSeq() - 1,
		}

		// add acked seq if synced
		if sequence.Synced() {
			resp.Ack = &storage.WriteResponse_AckSeq{AckSeq: sequence.GetAckSeq()}
		}

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func getLogicNodeFromCtx(ctx context.Context) (*models.Node, error) {
	return rpc.GetLogicNodeFromContext(ctx)
}

func parseCtx(ctx context.Context) (database string, shardID int32, logicNode *models.Node, err error) {
	logicNode, err = rpc.GetLogicNodeFromContext(ctx)
	if err != nil {
		return "", 0, nil, err
	}

	database, err = rpc.GetDatabaseFromContext(ctx)
	if err != nil {
		return "", 0, nil, err
	}

	shardID, err = rpc.GetShardIDFromContext(ctx)
	if err != nil {
		return "", 0, nil, err
	}
	return database, shardID, logicNode, err
}

func (w *Writer) getSequenceFromCtx(ctx context.Context) (replication.Sequence, error) {
	database, shardID, logicNode, err := parseCtx(ctx)
	if err != nil {
		return nil, err
	}
	return w.getSequence(database, shardID, *logicNode)
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
