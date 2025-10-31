package output

import (
	"github.com/lindb/common/pkg/logger"
)

type ConsoleOutput struct {
	logger logger.Logger
}

func NewConsoleOutput() Listener {
	return &ConsoleOutput{
		logger: logger.GetLogger("Streaming", "Console"),
	}
}

func (output *ConsoleOutput) Receive(event any) {
	output.logger.Info("console output, receive event", logger.Any("event", event))
}
