package page

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestNewFactory(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		listDirFunc = fileutil.ListDir
		mapFileFunc = fileutil.RWMap
	}()
	// case 1: list page files err
	listDirFunc = func(path string) ([]string, error) {
		return nil, fmt.Errorf("err")
	}
	fct, err := NewFactory(testPath, 128)
	assert.Error(t, err)
	assert.Nil(t, fct)
	// case 2: list page files parse file sequence err
	listDirFunc = func(path string) ([]string, error) {
		return []string{"a.bat"}, nil
	}
	fct, err = NewFactory(testPath, 128)
	assert.Error(t, err)
	assert.Nil(t, fct)
	// case 3: create page err
	listDirFunc = func(path string) ([]string, error) {
		return []string{"10.bat"}, nil
	}
	mapFileFunc = func(filePath string, size int) ([]byte, error) {
		return nil, fmt.Errorf("err")
	}
	fct, err = NewFactory(testPath, 128)
	assert.Error(t, err)
	assert.Nil(t, fct)
	// case 4: reopen page file
	listDirFunc = func(path string) ([]string, error) {
		return []string{"10.bat"}, nil
	}
	mapFileFunc = fileutil.RWMap
	fct, err = NewFactory(testPath, 128)
	assert.NoError(t, err)
	assert.NotNil(t, fct)
	fct1 := fct.(*factory)
	page, ok := fct1.pages[10]
	assert.True(t, ok)
	assert.NotNil(t, page)
}

func TestFactory_AcquirePage(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		mkDirFunc = fileutil.MkDirIfNotExist
		mapFileFunc = fileutil.RWMap
	}()
	// case 1: new factory err
	mkDirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	fct, err := NewFactory(testPath, 128)
	assert.Error(t, err)
	assert.Nil(t, fct)

	mkDirFunc = fileutil.MkDirIfNotExist

	// case 2: new factory success
	fct, err = NewFactory(testPath, 128)
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
	mapFileFunc = func(filePath string, size int) ([]byte, error) {
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

	// case 6: release page after close
	err = fct.ReleasePage(0)
	assert.Equal(t, errFactoryClosed, err)
}

func TestFactory_Close(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()

	fct, err := NewFactory(testPath, 128)
	assert.NoError(t, err)

	page1 := NewMockMappedPage(ctrl)
	fct1 := fct.(*factory)
	fct1.pages[1] = page1
	fct1.pages[2] = page1

	page1.EXPECT().Close().Return(fmt.Errorf("err")).MaxTimes(2)
	err = fct.Close()
	assert.NoError(t, err)
}

func TestFactory_ReleasePage(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		removeFileFunc = fileutil.RemoveFile
		ctrl.Finish()
	}()

	fct, err := NewFactory(testPath, 128)
	assert.NoError(t, err)
	p, err := fct.AcquirePage(10)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	files, err := fileutil.ListDir(testPath)
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	assert.Equal(t, int64(128), fct.Size())

	// remove file err
	removeFileFunc = func(file string) error {
		return fmt.Errorf("err")
	}
	err = fct.ReleasePage(10)
	assert.Error(t, err)
	files, err = fileutil.ListDir(testPath)
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	// remove file success
	removeFileFunc = fileutil.RemoveFile
	err = fct.ReleasePage(10)
	assert.NoError(t, err)
	files, err = fileutil.ListDir(testPath)
	assert.NoError(t, err)
	assert.Len(t, files, 0)

	assert.Equal(t, int64(0), fct.Size())

	err = fct.Close()
	assert.NoError(t, err)
}
