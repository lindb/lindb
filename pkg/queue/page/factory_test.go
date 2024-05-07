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

package page

import (
	"fmt"
	"os"
	"testing"

	commonfileutil "github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestNewFactory(t *testing.T) {
	tmpDir := t.TempDir()
	defer func() {
		listDirFunc = commonfileutil.ListDir
		mapFileFunc = fileutil.RWMap
	}()
	// case 1: list page files err
	listDirFunc = func(path string) ([]string, error) {
		return nil, fmt.Errorf("err")
	}
	fct, err := NewFactory(tmpDir, 128)
	assert.Error(t, err)
	assert.Nil(t, fct)
	// case 2: list page files parse file sequence err
	listDirFunc = func(path string) ([]string, error) {
		return []string{"a.bat"}, nil
	}
	fct, err = NewFactory(tmpDir, 128)
	assert.Error(t, err)
	assert.Nil(t, fct)
	// case 3: create page err
	listDirFunc = func(path string) ([]string, error) {
		return []string{"10.bat"}, nil
	}
	mapFileFunc = func(file *os.File, size int) ([]byte, error) {
		return nil, fmt.Errorf("err")
	}
	fct, err = NewFactory(tmpDir, 128)
	assert.Error(t, err)
	assert.Nil(t, fct)
	// case 4: reopen page file
	listDirFunc = func(path string) ([]string, error) {
		return []string{"10.bat"}, nil
	}
	mapFileFunc = fileutil.RWMap
	fct, err = NewFactory(tmpDir, 128)
	assert.NoError(t, err)
	assert.NotNil(t, fct)
	fct1 := fct.(*factory)
	page, ok := fct1.pages[10]
	assert.True(t, ok)
	assert.NotNil(t, page)
	err = fct.Close()
	assert.NoError(t, err)
}

func TestFactory_AcquirePage(t *testing.T) {
	tmpDir := t.TempDir()
	defer func() {
		mkDirFunc = commonfileutil.MkDirIfNotExist
		mapFileFunc = fileutil.RWMap
	}()
	// case 1: new factory err
	mkDirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	fct, err := NewFactory(tmpDir, 128)
	assert.Error(t, err)
	assert.Nil(t, fct)

	mkDirFunc = commonfileutil.MkDirIfNotExist

	// case 2: new factory success
	fct, err = NewFactory(tmpDir, 128)
	assert.NoError(t, err)
	assert.NotNil(t, fct)
	// case 3: acquire page success
	page1, err := fct.AcquirePage(0)
	assert.NoError(t, err)
	assert.NotNil(t, page1)
	p1, ok := fct.GetPage(0)
	assert.True(t, ok)
	assert.Equal(t, p1, page1)
	p1, ok = fct.GetPage(10)
	assert.False(t, ok)
	assert.Nil(t, p1)
	// get duplicate page
	page2, err := fct.AcquirePage(0)
	assert.NoError(t, err)
	assert.Equal(t, page1, page2)
	// case 4: get page err
	mapFileFunc = func(file *os.File, size int) ([]byte, error) {
		return nil, fmt.Errorf("err")
	}
	page2, err = fct.AcquirePage(2)
	assert.Error(t, err)
	assert.Nil(t, page2)
	mapFileFunc = fileutil.RWMap

	assert.Equal(t, int64(128), fct.Size())

	err = fct.Close()
	assert.NoError(t, err)
	// case 5: acquire page after close
	page2, err = fct.AcquirePage(2)
	assert.Equal(t, errFactoryClosed, err)
	assert.Nil(t, page2)
}

func TestFactory_Close(t *testing.T) {
	tmpDir := t.TempDir()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	fct, err := NewFactory(tmpDir, 128)
	assert.NoError(t, err)

	page1 := NewMockMappedPage(ctrl)
	fct1 := fct.(*factory)
	fct1.pages[1] = page1
	fct1.pages[2] = page1

	page1.EXPECT().Close().Return(fmt.Errorf("err")).MaxTimes(2)
	err = fct.Close()
	assert.NoError(t, err)
}

func TestFactory_TruncatePages(t *testing.T) {
	tmpDir := t.TempDir()
	ctrl := gomock.NewController(t)

	defer func() {
		removeFileFunc = commonfileutil.RemoveFile
		ctrl.Finish()
	}()

	fct, err := NewFactory(tmpDir, 128)
	assert.NoError(t, err)
	p, err := fct.AcquirePage(10)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	files, err := commonfileutil.ListDir(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	assert.Equal(t, int64(128), fct.Size())

	// remove file err
	removeFileFunc = func(file string) error {
		return fmt.Errorf("err")
	}
	fct.TruncatePages(11)
	files, err = commonfileutil.ListDir(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	// remove file success
	removeFileFunc = commonfileutil.RemoveFile
	fct.TruncatePages(11)
	assert.NoError(t, err)
	files, err = commonfileutil.ListDir(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, files, 0)

	assert.Equal(t, int64(0), fct.Size())

	err = fct.Close()
	assert.NoError(t, err)

	// truncate after closed
	fct.TruncatePages(11)

	page := NewMockMappedPage(ctrl)
	fct = &factory{
		pages:  map[int64]MappedPage{10: page},
		logger: logger.GetLogger("Queue", "Test"),
	}
	page.EXPECT().Close().Return(fmt.Errorf("err"))
	fct.TruncatePages(11)
}
