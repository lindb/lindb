package handler

import (
	"context"

	"github.com/eleme/lindb/replication"
	"github.com/eleme/lindb/rpc/proto/broker"
)

type Writer struct {
	cm replication.ChannelManager
}

func NewWriter(cm replication.ChannelManager) *Writer {
	return &Writer{
		cm: cm,
	}
}

func (w *Writer) Write(cxt context.Context, req *broker.WriteRequest) (*broker.WriteResponse, error) {
	ch, err := w.cm.GetChannel(req.GetCluster(), req.GetDatabase(), req.GetHash())
	if err != nil {
		return nil, err
	}

	if err := ch.Write(cxt, req.GetData()); err != nil {
		return nil, err
	}

	return &broker.WriteResponse{Msg: "Ok"}, nil
}
