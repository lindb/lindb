package query

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutorFactory_NewExecutor(t *testing.T) {
	factory := NewExecutorFactory()
	assert.NotNil(t, factory.NewStorageExecutor(context.TODO(), nil, nil, nil))
	assert.NotNil(t, factory.NewBrokerExecutor(context.TODO(), "db", "sql", nil, nil, nil))
}
