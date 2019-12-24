package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/tsdb"
)

func TestExecutorFactory_NewExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := NewExecutorFactory()
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	assert.NotNil(t, factory.NewStorageExecutor(nil, mockDatabase, nil, nil))
	assert.NotNil(t, factory.NewBrokerExecutor(
		context.TODO(), "db", "sql", nil, nil, nil))
	assert.NotNil(t, factory.NewMetadataStorageExecutor(nil, nil, nil))
	assert.NotNil(t, factory.NewMetadataBrokerExecutor(
		context.TODO(), "db", nil, nil, nil, nil))
}
