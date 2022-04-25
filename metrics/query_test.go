package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryStatistics(t *testing.T) {
	assert.NotNil(t, NewBrokerQueryStatistics())
	assert.NotNil(t, NewStorageQueryStatistics())
}
