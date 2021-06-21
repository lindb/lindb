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

package queue

import (
	"fmt"
	"math/rand"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/queue/page"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func randomString(length int) string {
	bytes := make([]byte, length)
	l := len(chars)

	for i := range bytes {
		bytes[i] = chars[rand.Intn(l)]
	}

	return string(bytes)
}

func TestFanOutQueue_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "fanOut")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		newQueueFunc = NewQueue
		mkDirFunc = fileutil.MkDirIfNotExist
		listDirFunc = fileutil.ListDir
		newFanOutFunc = NewFanOut

		ctrl.Finish()
	}()

	// case 1: create underlying queue err
	newQueueFunc = func(dirPath string, dataSizeLimit int64, removeTaskInterval time.Duration) (Queue, error) {
		return nil, fmt.Errorf("err")
	}
	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, fq)

	newQueueFunc = NewQueue
	// case 2: create fanOut path err
	queue := NewMockQueue(ctrl)
	queue.EXPECT().Close().AnyTimes()

	newQueueFunc = func(dirPath string, dataSizeLimit int64, removeTaskInterval time.Duration) (Queue, error) {
		return queue, nil
	}
	mkDirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	fq, err = NewFanOutQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, fq)

	mkDirFunc = fileutil.MkDirIfNotExist
	// case 3: list fanOut consumer group err
	listDirFunc = func(path string) ([]string, error) {
		return nil, fmt.Errorf("err")
	}
	fq, err = NewFanOutQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, fq)

	listDirFunc = fileutil.ListDir

	// case 4: create success
	queue.EXPECT().TailSeq().Return(int64(0))

	fq, err = NewFanOutQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	assert.NotNil(t, fq)
	fo, err := fq.GetOrCreateFanOut("group-1")
	assert.NoError(t, err)
	assert.NotNil(t, fo)
	fq.Close()
	// case 5: init fanOut consumer group err
	newFanOutFunc = func(parent, path string, q FanOutQueue) (FanOut, error) {
		return nil, fmt.Errorf("err")
	}
	fq, err = NewFanOutQueue(dir, 1024, time.Minute)
	assert.Error(t, err)
	assert.Nil(t, fq)
}

func TestFanOutQueue_GetOrCreateFanOut(t *testing.T) {
	dir := path.Join(testPath, "fanOut")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		newFanOutFunc = NewFanOut
	}()

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	assert.NotNil(t, fq)
	// case 1: create consumer group err
	newFanOutFunc = func(parent, path string, q FanOutQueue) (FanOut, error) {
		return nil, fmt.Errorf("err")
	}
	fo, err := fq.GetOrCreateFanOut("group-1")
	assert.Error(t, err)
	assert.Nil(t, fo)

	newFanOutFunc = NewFanOut

	// case 2: create consumer group success
	fo, err = fq.GetOrCreateFanOut("group-1")
	assert.NoError(t, err)
	assert.NotNil(t, fo)

	foNames := fq.FanOutNames()
	assert.Equal(t, "group-1", foNames[0])
}

func TestFanOutQueue_Sync(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "fanOut")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		err := fq.Put([]byte("12345"))
		assert.NoError(t, err)
	}

	assert.Equal(t, int64(-1), fq.TailSeq())
	// case 1: sync with empty consume group
	fq.Sync()
	assert.Equal(t, int64(-1), fq.TailSeq())
	fo1, err := fq.GetOrCreateFanOut("group-1")
	assert.NoError(t, err)
	fo1.Consume()       //0
	fo1.Consume()       //1
	s1 := fo1.Consume() //2
	fo2, err := fq.GetOrCreateFanOut("group-2")
	assert.NoError(t, err)

	s2 := fo2.Consume() //0

	// case 2: sync and ack min consume sequence, but consumer not ack
	fq.Sync()
	assert.Equal(t, int64(-1), fq.TailSeq())
	// case 3: ack and sync
	fo1.Ack(s1)
	fo2.Ack(s2)
	fq.Sync()
	assert.Equal(t, int64(0), fq.TailSeq())
	// case 3: sync err
	s2 = fo2.Consume() //1
	fo3 := fo2.(*fanOut)
	metaPage := page.NewMockMappedPage(ctrl)
	fo3.metaPage = metaPage
	metaPage.EXPECT().PutUint64(gomock.Any(), gomock.Any()).MaxTimes(2)
	metaPage.EXPECT().Sync().Return(fmt.Errorf("err"))
	fo2.Ack(s2)
}

func TestFanOut_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "fanOut")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	fo, err := fq.GetOrCreateFanOut("f1")
	assert.NoError(t, err)
	assert.NotNil(t, fo)

	fo1 := fo.(*fanOut)
	pageFct := page.NewMockFactory(ctrl)
	fo1.metaPageFct = pageFct
	pageFct.EXPECT().Close().Return(fmt.Errorf("err"))

	fq.Close()
}

func TestFanOut_new_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := path.Join(testPath, "fanOut")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		newPageFactoryFunc = page.NewFactory
		ctrl.Finish()
	}()

	// case 1: new meta page factory
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return nil, fmt.Errorf("err")
	}
	fo, err := NewFanOut(dir, "f1", nil)
	assert.Error(t, err)
	assert.Nil(t, fo)
	// case 2: acquire meta page err
	pageFct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return pageFct, nil
	}
	pageFct.EXPECT().Close().Return(fmt.Errorf("err"))
	pageFct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	fo, err = NewFanOut(dir, "f1", nil)
	assert.Error(t, err)
	assert.Nil(t, fo)
}

func TestFanOutQueue_one_consumer(t *testing.T) {
	dir := path.Join(testPath, "fanOut")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	assert.Empty(t, fq.FanOutNames())
	assert.Equal(t, int64(0), fq.HeadSeq())
	assert.Equal(t, int64(-1), fq.TailSeq())

	f1, err := fq.GetOrCreateFanOut("f1")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, fanOutDirName, "f1"), f1.Name())
	assert.Equal(t, int64(0), f1.HeadSeq())
	assert.Equal(t, int64(-1), f1.TailSeq())
	assert.Equal(t, SeqNoNewMessageAvailable, f1.Consume())
	assert.Equal(t, int64(0), f1.Pending())

	// msg 0
	msg := []byte("123")
	err = fq.Put(msg)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), f1.Pending())

	fseq := f1.Consume()
	assert.Equal(t, int64(0), fseq)
	assert.Equal(t, int64(1), f1.HeadSeq())
	assert.Equal(t, int64(-1), f1.TailSeq())
	assert.Equal(t, int64(0), f1.Pending())

	fmsg, err := f1.Get(0)
	assert.NoError(t, err)
	assert.Equal(t, msg, fmsg)
	assert.Equal(t, SeqNoNewMessageAvailable, f1.Consume())

	// msg1, msg2
	msg1 := []byte("456")
	msg2 := []byte("789")

	err = fq.Put(msg1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), fq.HeadSeq())
	assert.Equal(t, int64(1), f1.Pending())

	err = fq.Put(msg2)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), fq.HeadSeq())
	assert.Equal(t, int64(2), f1.Pending())

	fseq = f1.Consume()
	assert.Equal(t, int64(1), fseq)
	assert.Equal(t, int64(2), f1.HeadSeq())
	assert.Equal(t, int64(1), f1.Pending())

	fmsg, err = f1.Get(fseq)
	assert.NoError(t, err)
	assert.Equal(t, msg1, fmsg)

	f1.Ack(fseq) // ack 1
	assert.Equal(t, fseq, f1.TailSeq())

	fseq = f1.Consume()
	assert.Equal(t, int64(2), fseq)
	assert.Equal(t, int64(3), f1.HeadSeq())
	assert.Equal(t, int64(0), f1.Pending())

	fmsg, err = f1.Get(fseq)
	assert.NoError(t, err)
	assert.Equal(t, msg2, fmsg)
	f1.Ack(fseq) // akc 2
	assert.Equal(t, fseq, f1.TailSeq())
	assert.Equal(t, int64(0), f1.Pending())

	fq.Close()
	// reopen
	fq, err = NewFanOutQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	f1, err = fq.GetOrCreateFanOut("f1")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), f1.TailSeq())
	assert.Equal(t, int64(3), f1.HeadSeq())
	assert.Equal(t, int64(0), f1.Pending())
	fq.Close()
}

func TestFanOutQueue_SetHeadSeq(t *testing.T) {
	dir := path.Join(testPath, "fanOut")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)

	f1, err := fq.GetOrCreateFanOut("f1")
	assert.NoError(t, err)

	err = f1.SetHeadSeq(1)
	assert.Error(t, err)

	err = fq.Put([]byte("123"))
	assert.NoError(t, err)

	err = fq.Put([]byte("456"))
	assert.NoError(t, err)

	seq := f1.Consume()
	assert.Equal(t, int64(0), seq)

	seq = f1.Consume()
	assert.Equal(t, int64(1), seq)

	// reset head consume sequence
	err = f1.SetHeadSeq(-1)
	assert.NoError(t, err)

	seq = f1.Consume()
	assert.Equal(t, int64(0), seq)

	seq = f1.Consume()
	assert.Equal(t, int64(1), seq)

	f1.Ack(1)

	err = f1.SetHeadSeq(0)
	assert.Error(t, err)
	fq.Close()
}

func TestFanOutQueue_multiple_consumer(t *testing.T) {
	dir := path.Join(testPath, "fanOut")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)
	// put data
	for i := 0; i < 100; i++ {
		msg := []byte(fmt.Sprintf("msg-%d", i))
		err = fq.Put(msg)
		assert.NoError(t, err)
	}

	// consumer group 1
	consumeMsg(t, fq, "f1", 100)
	// consumer group 2
	consumeMsg(t, fq, "f2", 100)

	fq.Close()
}

func TestFanOutQueue_SetAppendSeq(t *testing.T) {
	dir := path.Join(testPath, "fanOut")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	assert.NoError(t, err)

	// put data
	for i := 0; i < 100; i++ {
		msg := []byte(fmt.Sprintf("msg-%d", i))
		err = fq.Put(msg)
		assert.NoError(t, err)
	}

	// consumer group 1
	consumeMsg(t, fq, "f1", 100)
	f2, _ := fq.GetOrCreateFanOut("f2")
	// set new append seq
	fq.SetAppendSeq(200)
	// put msg
	msg := []byte(fmt.Sprintf("msg-%d", 200))
	err = fq.Put(msg)
	assert.NoError(t, err)

	f1, _ := fq.GetOrCreateFanOut("f1")
	fseq := f1.Consume()
	assert.Equal(t, fseq, int64(200))
	fmsg, err := f1.Get(fseq)
	assert.NoError(t, err)
	assert.Equal(t, []byte(fmt.Sprintf("msg-%d", 200)), fmsg)

	fseq = f2.Consume()
	assert.Equal(t, fseq, int64(200))
	fmsg, err = f2.Get(fseq)
	assert.NoError(t, err)
	assert.Equal(t, []byte(fmt.Sprintf("msg-%d", 200)), fmsg)

	fq.Close()
}

func TestFanOutQueue_concurrent_read(t *testing.T) {
	dir := path.Join(testPath, "fanout_concurrent")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	msgSize := 1024
	dataFileSize := int64(512)
	// random text
	bytesSli := make([][]byte, msgSize)

	for i := range bytesSli {
		bytesSli[i] = []byte(randomString(rand.Intn(10) + 1))
	}

	fq, err := NewFanOutQueue(dir, dataFileSize, time.Second)
	assert.NoError(t, err)

	readConcurrent := 10
	wg := &sync.WaitGroup{}
	wg.Add(readConcurrent)

	for i := 0; i < readConcurrent; i++ {
		go func(seq int) {
			defer wg.Done()
			read(t, bytesSli, fq, "fo-"+strconv.Itoa(seq))
		}(i)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		write(t, bytesSli, fq)
	}()

	wg.Wait()

	assert.Equal(t, int64(msgSize), fq.HeadSeq())

	// wait for background deleting
	time.Sleep(2 * time.Second)

	fq.Close()

	// reload
	fq2, err := NewFanOutQueue(dir, dataFileSize, time.Second)
	assert.NoError(t, err)

	assert.Equal(t, int64(msgSize), fq2.HeadSeq())

	assert.Equal(t, len(fq2.FanOutNames()), readConcurrent)

	for i := 0; i < readConcurrent; i++ {
		fo, err := fq2.GetOrCreateFanOut("fo-" + strconv.Itoa(i))
		assert.NoError(t, err)

		assert.Equal(t, int64(msgSize), fo.HeadSeq())
		assert.Equal(t, int64(msgSize)-1, fo.TailSeq())
	}
	fq2.Close()
}

func read(t *testing.T, raw [][]byte, fq FanOutQueue, name string) {
	fo, err := fq.GetOrCreateFanOut(name)
	assert.NoError(t, err)

	counter := 0

	for {
		if counter == len(raw) {
			return
		}

		seq := fo.Consume()

		if seq == SeqNoNewMessageAvailable {
			time.Sleep(time.Microsecond * 10)
			continue
		}

		assert.Equal(t, seq, int64(counter))

		bys, err := fo.Get(seq)
		assert.NoError(t, err)
		assert.Equal(t, raw[int(seq)], bys)

		counter++

		fo.Ack(seq)
	}
}

func write(t *testing.T, raw [][]byte, fq FanOutQueue) {
	for _, bys := range raw {
		err := fq.Put(bys)
		assert.NoError(t, err)
	}
}

func consumeMsg(t *testing.T, fq FanOutQueue, consumerGroup string, msgCount int) {
	f1, err := fq.GetOrCreateFanOut(consumerGroup)
	assert.NoError(t, err)

	for i := 0; i < msgCount; i++ {
		fseq := f1.Consume()
		assert.Equal(t, fseq, int64(i))
		fmsg, err := f1.Get(int64(i))
		assert.NoError(t, err)
		assert.Equal(t, []byte(fmt.Sprintf("msg-%d", i)), fmsg)
	}
}
