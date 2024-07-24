package execution

import (
	"fmt"
	"sync"
	"time"

	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/execution/model"
)

type RequestIDGenerator struct {
	node    string
	seq     int
	current int64
	next    int64

	lock sync.Mutex
}

func NewRequestIDGenerator(node string) *RequestIDGenerator {
	return &RequestIDGenerator{
		node: node,
		seq:  0,
	}
}

func (g *RequestIDGenerator) GenerateRequestID() model.RequestID {
	now := time.Now().UnixMilli()

	g.lock.Lock()
	defer g.lock.Unlock()

	if now >= g.next {
		// initialize next loop
		g.initialize(now)
	}

	g.seq++
	return model.RequestID(
		fmt.Sprintf("%s-%s-%08d",
			g.node,
			timeutil.FormatTimestamp(g.current, timeutil.DataTimeFormat4),
			g.seq,
		))
}

func (g *RequestIDGenerator) initialize(now int64) {
	g.current = now - now%timeutil.OneSecond
	g.next = g.current + timeutil.OneSecond
	g.seq = 0
}
