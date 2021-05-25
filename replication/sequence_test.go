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

package replication

import (
	"fmt"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/queue/page"
)

var testPath = "test"

func TestSequence_new_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	tmp := path.Join(testPath, "sequence_test")
	defer func() {
		newPageFactoryFunc = page.NewFactory
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()
	// case 1: new page factory err
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return nil, fmt.Errorf("err")
	}
	seq, err := NewSequence(tmp)
	assert.Error(t, err)
	assert.Nil(t, seq)
	// case 2: AcquirePage err
	fct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return fct, nil
	}
	fct.EXPECT().GetPage(int64(metaPageID)).Return(nil, false)
	fct.EXPECT().Close().Return(fmt.Errorf("err"))
	fct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	seq, err = NewSequence(tmp)
	assert.Error(t, err)
	assert.Nil(t, seq)
	// case 3: sync err
	fct.EXPECT().GetPage(int64(metaPageID)).Return(nil, false)
	fct.EXPECT().Close().Return(fmt.Errorf("err"))
	mockPage := page.NewMockMappedPage(ctrl)
	mockPage.EXPECT().PutUint64(gomock.Any(), gomock.Any())
	mockPage.EXPECT().Sync().Return(fmt.Errorf("err"))
	fct.EXPECT().AcquirePage(gomock.Any()).Return(mockPage, nil)
	seq, err = NewSequence(tmp)
	assert.Error(t, err)
	assert.Nil(t, seq)
}

func TestSequence(t *testing.T) {
	tmp := path.Join(testPath, "sequence_test")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	seq, err := NewSequence(tmp)
	assert.NoError(t, err)
	assert.NotNil(t, seq)

	assert.Equal(t, seq.GetHeadSeq(), int64(-1))
	assert.Equal(t, seq.GetAckSeq(), int64(-1))
	err = seq.Close()
	assert.NoError(t, err)
	seq, err = NewSequence(tmp)
	assert.NoError(t, err)
	assert.NotNil(t, seq)

	assert.Equal(t, seq.GetHeadSeq(), int64(-1))
	assert.Equal(t, seq.GetAckSeq(), int64(-1))

	seq.SetHeadSeq(int64(10))
	seq.SetAckSeq(int64(5))

	assert.Equal(t, seq.GetHeadSeq(), int64(10))
	assert.Equal(t, seq.GetAckSeq(), int64(5))

	err = seq.Sync()
	assert.NoError(t, err)

	// new sequence
	newSeq, err := NewSequence(tmp)
	assert.NoError(t, err)

	assert.Equal(t, newSeq.GetAckSeq(), int64(5))
	assert.Equal(t, newSeq.GetHeadSeq(), int64(5))
}
