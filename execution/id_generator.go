package execution

import (
	"fmt"
	"time"

	"github.com/lindb/common/pkg/timeutil"
)

type RequestID string

type RequestIDGenerator struct {
	node string
}

func NewRequestIDGenerator(node string) *RequestIDGenerator {
	return &RequestIDGenerator{
		node: node,
	}
}

func (g RequestIDGenerator) GenerateRequestID() RequestID {
	now := time.Now().UnixMilli()
	return RequestID(
		fmt.Sprintf("%s-%05d-%s",
			timeutil.FormatTimestamp(now, timeutil.DataTimeFormat4),
			1,
			g.node,
		))
}
