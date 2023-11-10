// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package metadb

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewMetadata(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer func() {
		createMetadataBackendFn = newMetadataBackend
		ctrl.Finish()
	}()
	metadata1, err := NewMetadata(context.TODO(), "test", testPath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, metadata1.TagMetadata())
	assert.NotNil(t, metadata1.MetadataDatabase())
	assert.Equal(t, "test", metadata1.DatabaseName())

	createMetadataBackendFn = func(parent string) (MetadataBackend, error) {
		return nil, fmt.Errorf("err")
	}
	metadata2, err := NewMetadata(context.TODO(), "test", testPath, nil)
	assert.Error(t, err)
	assert.Nil(t, metadata2)

	err = metadata1.Close()
	assert.NoError(t, err)

	db := NewMockMetadataDatabase(ctrl)
	m := metadata1.(*metadata)
	m.metadataDatabase = db
	db.EXPECT().Close().Return(fmt.Errorf("err"))
	err = m.Close()
	assert.Error(t, err)
}

func TestMetadata_Flush(t *testing.T) {
	testPath := t.TempDir()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metadata1, err := NewMetadata(context.TODO(), "test", testPath, nil)
	assert.NoError(t, err)
	db := NewMockMetadataDatabase(ctrl)
	m := metadata1.(*metadata)
	backendDB := m.metadataDatabase
	defer func() {
		_ = backendDB.Close()
	}()
	m.metadataDatabase = db
	db.EXPECT().Sync().Return(fmt.Errorf("err"))
	err = metadata1.Flush()
	assert.Error(t, err)

	db.EXPECT().Sync().Return(nil)
	err = metadata1.Flush()
	assert.NoError(t, err)
}
