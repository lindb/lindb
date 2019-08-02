package handler

import (
	"context"
	"fmt"
	"io"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
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
		logger:         logger.GetLogger("handler/writer"),
	}
}

// Write handles the stream write request.
func (w *Writer) Write(stream storage.WriteService_WriteServer) error {

	database, shardID, logicNode, err := parseContext(stream.Context())
	if err != nil {
		return err
	}

	sequence, ok := w.sm.GetSequence(database, shardID, *logicNode)
	if !ok {
		var err error
		sequence, err = w.sm.CreateSequence(database, shardID, *logicNode)
		if err != nil {
			return err
		}
	}

	shard := w.storageService.GetShard(database, shardID)
	if shard == nil {
		return fmt.Errorf("shard %d for database %s not exists", shardID, database)
	}

	// only send resetSeq response once for a sequenceNum.
	var resetSeq int64 = -1
	var lastResetSeq int64 = -1
	for {
		req, err := stream.Recv()
		w.logger.Info("req", logger.Any("req", req))
		if err == io.EOF {
			return nil
		}
		if err != nil {
			w.logger.Error("error", logger.Error(err))
			return err
		}

		for _, rep := range req.Replicas {
			seq := rep.Seq
			if sequence.GetHeadSeq() == 0 {
				sequence.SetHeadSeq(seq)
			}

			hs := sequence.GetHeadSeq()
			if hs != seq {
				// reset to headSeq
				resetSeq = hs
				break
			}

			sequence.SetHeadSeq(hs + 1)

			// todo write to storage
			//data := rep.Data
		}

		// reset seq, only one reset response for a seqNum.
		if resetSeq != -1 && resetSeq != lastResetSeq {
			err := stream.Send(&storage.WriteResponse{
				Seq: &storage.WriteResponse_ResetSeq{
					ResetSeq: resetSeq,
				},
			})

			if err != nil {
				return err
			}

			lastResetSeq = resetSeq
		}

		// acked seq
		if sequence.Synced() {
			err := stream.Send(&storage.WriteResponse{
				Seq: &storage.WriteResponse_AckSeq{
					AckSeq: sequence.GetAckSeq(),
				},
			})
			if err != nil {
				return err
			}
		}
	}
}

func parseContext(ctx context.Context) (database string, shardID int32, logicNode *models.Node, err error) {
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
