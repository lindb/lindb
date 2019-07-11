package handler

import (
	"context"

	"github.com/eleme/lindb/rpc"
	"github.com/eleme/lindb/rpc/proto/common"
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

func (w *Writer) WritePoints(ctx context.Context, request *common.Request) (*common.Response, error) {
	// todo: @XiaTianliang
	//bs.logger.Info(string(request.Data))
	return rpc.ResponseOK(), nil
}
