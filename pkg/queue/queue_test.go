package queue

import (
	"fmt"
	"path"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/queue/page"
)

var testPath = "test"

func TestQueue_Put(t *testing.T) {
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	q, err := NewQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	// case 1: init queue
	assert.Equal(t, int64(0), q.Size())
	assert.Equal(t, int64(-1), q.HeadSeq())
	assert.Equal(t, int64(-1), q.TailSeq())
	// case 2: put data
	err = q.Put([]byte("123"))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), q.Size())
	assert.Equal(t, int64(0), q.HeadSeq())
	assert.Equal(t, int64(-1), q.TailSeq())
	// case 4: put data
	err = q.Put([]byte("456"))
	assert.NoError(t, err)
	assert.Equal(t, int64(2), q.Size())
	assert.Equal(t, int64(1), q.HeadSeq())
	assert.Equal(t, int64(-1), q.TailSeq())
	// read data
	data, err := q.Get(0)
	assert.NoError(t, err)
	assert.Equal(t, []byte("123"), data)
	data, err = q.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, []byte("456"), data)
	data, err = q.Get(2)
	assert.Error(t, err)
	assert.Nil(t, data)
	q.Close()
	// case 5: re-open
	q, err = NewQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	err = q.Put([]byte("789"))
	assert.NoError(t, err)
	assert.Equal(t, int64(3), q.Size())
	assert.Equal(t, int64(2), q.HeadSeq())
	assert.Equal(t, int64(-1), q.TailSeq())
	// case 6: get message
	data, err = q.Get(0)
	assert.NoError(t, err)
	assert.Equal(t, []byte("123"), data)
	data, err = q.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, []byte("456"), data)
	data, err = q.Get(2)
	assert.NoError(t, err)
	assert.Equal(t, []byte("789"), data)
	q.Close()
	// case 6: re-open can read data
	q, err = NewQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	data, err = q.Get(0)
	assert.NoError(t, err)
	assert.Equal(t, []byte("123"), data)
	data, err = q.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, []byte("456"), data)
	data, err = q.Get(2)
	assert.NoError(t, err)
	assert.Equal(t, []byte("789"), data)
	assert.False(t, q.IsEmpty())
	q.Close()
}

func TestQueue_Ack(t *testing.T) {
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	q, err := NewQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	err = q.Put([]byte("123"))
	assert.NoError(t, err)
	err = q.Put([]byte("456"))
	assert.NoError(t, err)
	data, err := q.Get(0)
	assert.NoError(t, err)
	assert.Equal(t, []byte("123"), data)
	data, err = q.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, []byte("456"), data)
	q.Ack(1)
	assert.Equal(t, int64(0), q.Size())
	assert.True(t, q.IsEmpty())

	for i := 0; i < 10; i++ {
		data, err = q.Get(int64(i))
		assert.Equal(t, ErrOutOfSequenceRange, err)
		assert.Nil(t, data)
	}
	q.Close()
}

func TestQueue_new_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		mkDirFunc = fileutil.MkDirIfNotExist
		newPageFactoryFunc = page.NewFactory

		ctrl.Finish()
	}()

	// case 1: create path fail
	mkDirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	q, err := NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)

	mkDirFunc = fileutil.MkDirIfNotExist
	// case 2: create data page factory err
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return nil, fmt.Errorf("err")
	}
	q, err = NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)
	// case 3: create index page factory err
	fct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		if strings.HasSuffix(path, dataPath) {
			return fct, nil
		}

		return nil, fmt.Errorf("err")
	}

	fct.EXPECT().Close().Return(nil)

	q, err = NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)
	// case 4: create meta page factory err
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		if strings.HasSuffix(path, dataPath) {
			return fct, nil
		} else if strings.HasSuffix(path, indexPath) {
			return fct, nil
		}

		return nil, fmt.Errorf("err")
	}

	fct.EXPECT().Close().Return(nil).MaxTimes(2)

	q, err = NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)
	// case 5: acquire meta page err
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return fct, nil
	}

	fct.EXPECT().Close().Return(nil).MaxTimes(3)
	fct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))

	q, err = NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)
	// case 6: acquire index page err when empty queue
	indexFct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		if strings.HasSuffix(path, indexPath) {
			return indexFct, nil
		}

		return page.NewFactory(path, pageSize)
	}

	indexFct.EXPECT().Close().Return(nil)
	indexFct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))

	q, err = NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)
	// case 6: acquire data page err when empty queue
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		if strings.HasSuffix(path, dataPath) {
			return fct, nil
		}

		return page.NewFactory(path, pageSize)
	}

	fct.EXPECT().Close().Return(nil)
	fct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))

	q, err = NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)
	// case 7: sync meta data err
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		if strings.HasSuffix(path, metaPath) {
			return fct, nil
		}

		return page.NewFactory(path, pageSize)
	}

	fct.EXPECT().Close().Return(nil)

	metaPage := page.NewMockMappedPage(ctrl)

	fct.EXPECT().AcquirePage(gomock.Any()).Return(metaPage, nil)
	metaPage.EXPECT().PutUint64(gomock.Any(), gomock.Any()).MaxTimes(4)
	metaPage.EXPECT().Sync().Return(fmt.Errorf("err"))
	// remove old data
	_ = fileutil.RemoveDir(testPath)
	q, err = NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)
}

func TestQueue_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	pageFct := page.NewMockFactory(ctrl)
	pageFct.EXPECT().Close().Return(fmt.Errorf("err")).MaxTimes(3)
	q, err := NewQueue(dir, 1024, time.Minute)
	q1 := q.(*queue)
	q1.dataPageFct = pageFct
	q1.indexPageFct = pageFct
	q1.metaPageFct = pageFct
	assert.NoError(t, err)
	q.Close()
}

func TestQueue_reopen_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		newPageFactoryFunc = page.NewFactory

		ctrl.Finish()
	}()

	q, err := NewQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	err = q.Put([]byte("123"))
	assert.NoError(t, err)
	q.Close()

	// case 1: acquire index page err
	fct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		if strings.HasSuffix(path, indexPath) {
			return fct, nil
		}

		return page.NewFactory(path, pageSize)
	}

	fct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	fct.EXPECT().Close().Return(nil)

	q, err = NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)

	// case 1: acquire data page err
	fct = page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		if strings.HasSuffix(path, dataPath) {
			return fct, nil
		}

		return page.NewFactory(path, pageSize)
	}

	fct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	fct.EXPECT().Close().Return(nil)

	q, err = NewQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, q)
}

func TestQueue_Put_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()

	mockPage := page.NewMockMappedPage(ctrl)
	mockPage.EXPECT().Sync().Return(fmt.Errorf("err")).AnyTimes()
	q, err := NewQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	// case 1: data > page size, return err
	data := make([]byte, dataPageSize+10)
	err = q.Put(data)
	assert.Error(t, err)
	// case 2: alloc new data page err
	err = q.Put([]byte("123456789"))
	assert.NoError(t, err)

	q1 := q.(*queue)
	q1.dataPage = mockPage
	q1.indexPage = mockPage
	fct := page.NewMockFactory(ctrl)
	fct.EXPECT().Size().Return(int64(1000)).AnyTimes()
	dataFct := q1.dataPageFct
	q1.dataPageFct = fct

	fct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))

	data = make([]byte, dataPageSize-5)
	err = q.Put(data)
	assert.Error(t, err)

	q1.dataPage = mockPage
	q1.dataPageFct = dataFct
	// case 3: alloc new index page err
	indexFct := page.NewMockFactory(ctrl)
	indexFct.EXPECT().Size().Return(int64(1000)).AnyTimes()
	indexFct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	q1.indexPageFct = indexFct
	q1.headSeq.Store(indexItemsPerPage)

	err = q.Put(data)
	assert.Error(t, err)

	indexFct.EXPECT().Close().Return(nil)
	q.Close()
}

func TestQueue_Get_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()

	q, err := NewQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	err = q.Put([]byte("123456789"))
	assert.NoError(t, err)

	fct := page.NewMockFactory(ctrl)
	q1 := q.(*queue)
	indexFct := q1.indexPageFct
	q1.indexPageFct = fct
	// case 1: index page not exist
	fct.EXPECT().GetPage(gomock.Any()).Return(nil, false)

	data, err := q.Get(0)
	assert.Error(t, err)
	assert.Nil(t, data)

	q1.indexPageFct = indexFct
	// case 2: data page not exist
	dataFct := q1.dataPageFct
	q1.dataPageFct = fct

	fct.EXPECT().GetPage(gomock.Any()).Return(nil, false)

	data, err = q.Get(0)
	assert.Error(t, err)
	assert.Nil(t, data)

	q1.dataPageFct = dataFct

	q.Close()
}

func TestQueue_Ack_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		newPageFactoryFunc = page.NewFactory

		ctrl.Finish()
	}()

	q, err := NewQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	err = q.Put([]byte("123456789"))
	assert.NoError(t, err)

	mockMetaPage := page.NewMockMappedPage(ctrl)
	q1 := q.(*queue)
	metaPage := q1.metaPage
	q1.metaPage = mockMetaPage

	// sync meta page err
	mockMetaPage.EXPECT().PutUint64(gomock.Any(), gomock.Any())
	mockMetaPage.EXPECT().Sync().Return(fmt.Errorf("err"))
	q.Ack(0)

	q1.metaPage = metaPage

	q.Close()
}

func TestQueue_data_limit(t *testing.T) {
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	q, err := NewQueue(dir, 128*1024*1024, time.Second)
	assert.NoError(t, err)
	q1 := q.(*queue)
	q1.dataSizeLimit = dataPageSize - 10
	data := make([]byte, dataPageSize-10)
	// put data
	err = q.Put(data)
	assert.NoError(t, err)
	// need acquire data page, but size limit
	err = q.Put(data)
	assert.Equal(t, ErrExceedingTotalSizeLimit, err)

	q1.dataSizeLimit = 2 * dataPageSize
	q1.headSeq.Store(indexItemsPerPage)
	// need acquire index page, but size limit
	err = q.Put(data)
	assert.Equal(t, ErrExceedingTotalSizeLimit, err)
}

func TestQueue_concurrently(t *testing.T) {
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	q, err := NewQueue(dir, 128*1024*1024, time.Second)
	assert.NoError(t, err)

	var (
		messages sync.Map
		wait     sync.WaitGroup
	)

	sendMessages := make([]map[string][]byte, 4)

	for i := 0; i < 4; i++ {
		sendMessages[i] = mockMessageData(i, 100)
	}
	wait.Add(4)

	for i := 0; i < 4; i++ {
		msg := sendMessages[i]

		go func() {
			defer wait.Done()

			for k, v := range msg {
				err := q.Put(v)
				messages.Store(k, v)

				if err != nil {
					panic("get err")
				}
			}
		}()
	}

	wait.Wait()
	assert.Equal(t, int64(400), q.Size())

	for i := 0; i < 400; i++ {
		data, err := q.Get(int64(i))
		assert.NoError(t, err)
		messages.Delete(string(data))
	}

	messages.Range(func(key, value interface{}) bool {
		panic("get data")
	})
}

func TestQueue_remove_expire_page(t *testing.T) {
	dir := path.Join(testPath, "queue")
	ctrl := gomock.NewController(t)

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()

	indexPageFct := page.NewMockFactory(ctrl)
	dataPageFct := page.NewMockFactory(ctrl)
	metaPage := page.NewMockMappedPage(ctrl)
	q, err := NewQueue(dir, dataPageSize*8, 500*time.Second)
	assert.NoError(t, err)
	q.Close()

	q1 := q.(*queue)
	q1.metaPage = metaPage
	q1.indexPageFct = indexPageFct
	q1.dataPageFct = dataPageFct
	// case 1: ack sequence < 0
	q1.removeExpirePage()
	q1.tailSeq.Store(indexItemsPerPage * 3)
	// case 2: index page not exist
	indexPageFct.EXPECT().GetPage(gomock.Any()).Return(nil, false)
	q1.removeExpirePage()
	// case 3: release page
	indexPage := page.NewMockMappedPage(ctrl)
	gomock.InOrder(
		indexPageFct.EXPECT().GetPage(gomock.Any()).Return(indexPage, true),
		indexPage.EXPECT().ReadUint64(gomock.Any()).Return(uint64(3)),
		dataPageFct.EXPECT().ReleasePage(int64(0)).Return(nil),
		metaPage.EXPECT().PutUint64(uint64(0), queueExpireDataOffset),
		dataPageFct.EXPECT().ReleasePage(int64(1)).Return(nil),
		metaPage.EXPECT().PutUint64(uint64(1), queueExpireDataOffset),
		dataPageFct.EXPECT().ReleasePage(int64(2)).Return(fmt.Errorf("err")),
		indexPageFct.EXPECT().ReleasePage(int64(0)).Return(nil),
		metaPage.EXPECT().PutUint64(uint64(0), queueExpireIndexOffset),
		indexPageFct.EXPECT().ReleasePage(int64(1)).Return(fmt.Errorf("err")),
		metaPage.EXPECT().Sync().Return(fmt.Errorf("err")),
	)
	q1.removeExpirePage()
}

func TestQueue_big_loop(t *testing.T) {
	dir := path.Join(testPath, "queue")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	q, err := NewQueue(dir, dataPageSize*8, 500*time.Millisecond)
	assert.NoError(t, err)
	loop := 1000000
	str := "big_loop_test"
	for i := 0; i < loop; i++ {
		err = q.Put([]byte(fmt.Sprintf("%s-%d", str, i)))
		assert.NoError(t, err)
	}
	for i := 0; i < loop; i++ {
		data, err := q.Get(int64(i))
		assert.NoError(t, err)
		assert.Equal(t, []byte(fmt.Sprintf("%s-%d", str, i)), data)
	}
	q.Ack(1000000 - 10)
	time.Sleep(time.Second)

	q.Close()
}

func mockMessageData(bucket, length int) map[string][]byte {
	data := make(map[string][]byte)

	for i := 0; i < length; i++ {
		str := fmt.Sprintf("%d-bucket-%d", bucket, i)
		data[str] = []byte(str)
	}

	return data
}
