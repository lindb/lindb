package monitoring

import (
	"testing"

	"github.com/lindb/lindb/internal/linmetric"
)

func TestRuntimeObserver_Observe(t *testing.T) {
	r := newRuntimeObserver(linmetric.BrokerRegistry)
	r.Observe()
}
