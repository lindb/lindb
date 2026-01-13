package cep

import (
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/streaming/cep/annotation"
)

type Listener struct {
	mapper *annotation.MetricMapper
	logger logger.Logger
}

func NewListener() *Listener {
	return &Listener{
		mapper: annotation.NewMetricMapper(),
		logger: logger.GetLogger("CEP", "Listener"),
	}
}

func (l *Listener) Receive(event any) {
	page, ok := event.(*types.Page)
	if !ok {
		return
	}
	l.mapper.Map(page)
	l.logger.Info("Listener, receive event", logger.Any("event", page))
}
