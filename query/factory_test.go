package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestExecutorFactory_NewExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := NewExecutorFactory()
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	assert.NotNil(t, factory.NewStorageExecutor(nil, mockDatabase, newStorageExecuteContext("ns", nil, &stmt.Query{})))
	assert.NotNil(t, factory.NewBrokerExecutor(
		context.TODO(), "db", "ns", "sql", nil, nil, nil, nil))
	assert.NotNil(t, factory.NewMetadataStorageExecutor(nil, "ns", nil, nil))
	assert.NotNil(t, factory.NewMetadataBrokerExecutor(
		context.TODO(), "db", "ns", nil, nil, nil, nil))
}

func TestNewExecutorFactory_NewContext(t *testing.T) {
	factory := NewExecutorFactory()
	assert.NotNil(t, factory.NewStorageExecuteContext("ns", nil, &stmt.Query{}))
}
