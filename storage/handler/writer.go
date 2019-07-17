package handler

import (
	"google.golang.org/grpc/peer"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/rpc/proto/storage"
	"github.com/eleme/lindb/service"
)

type Writer struct {
	storageService service.StorageService
}

func NewWriter(storageService service.StorageService) *Writer {
	return &Writer{
		storageService: storageService,
	}
}

func (w *Writer) Write(stream storage.WriteService_WriteServer) error {
	remotePeer, ok := peer.FromContext(stream.Context())
	log := logger.GetLogger("storage/writer")
	log.Info("peer", logger.Any("peer", remotePeer), logger.Any("ok", ok), logger.Any("ctx", stream.Context()))

	for {
		req, err := stream.Recv()
		log.Info("req", logger.Any("req", req))
		if err != nil {
			log.Error("error", logger.Error(err))
			return err
		}
		replica := req.Replica
		if replica != nil {
			database := w.storageService.GetEngine(replica.Database)
			if database != nil {
				shard := database.GetShard(replica.ShardID)
				if shard != nil {
					_ = shard.Write(nil)
				}
			}
		}
	}
}
