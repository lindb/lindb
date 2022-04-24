package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/linmetric"
)

func TestNewRuntimeStatistics(t *testing.T) {
	assert.NotNil(t, NewRuntimeStatistics(linmetric.BrokerRegistry))
}
