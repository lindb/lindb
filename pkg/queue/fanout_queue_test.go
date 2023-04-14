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
	dir := filepath.Join(t.TempDir(), t.Name())

	defer func() {
		newQueueFunc = NewQueue
		mkDirFunc = fileutil.MkDirIfNotExist
		listDirFunc = fileutil.ListDir
		newConsumerGroupFunc = NewConsumerGroup

		ctrl.Finish()
	}()

	// case 1: create underlying queue err
	newQueueFunc = func(dirPath string, dataSizeLimit int64) (Queue, error) {
		return nil, fmt.Errorf("err")
	}
	fq, err := NewFanOutQueue(dir, 1024)
	assert.Error(t, err)
	assert.Nil(t, fq)

	newQueueFunc = NewQueue
	// case 2: create consumerGroup path err
	queue := NewMockQueue(ctrl)
	queue.EXPECT().Close().AnyTimes()
	queue.EXPECT().Signal().AnyTimes()

	newQueueFunc = func(dirPath string, dataSizeLimit int64) (Queue, error) {
		return queue, nil
	}
	mkDirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	fq, err = NewFanOutQueue(dir, 1024)
	assert.Error(t, err)
	assert.Nil(t, fq)

	mkDirFunc = fileutil.MkDirIfNotExist
	// case 3: list consumerGroup consumer group err
	listDirFunc = func(path string) ([]string, error) {
		return nil, fmt.Errorf("err")
	}
	fq, err = NewFanOutQueue(dir, 1024)
	assert.Error(t, err)
	assert.Nil(t, fq)

	listDirFunc = fileutil.ListDir

	// case 4: create success
	fq, err = NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	assert.NotNil(t, fq)
	assert.Equal(t, dir, fq.Path())
	fo, err := fq.GetOrCreateConsumerGroup("group-1")
	assert.NoError(t, err)
	assert.NotNil(t, fo)
	fq.Close()
	// case 5: init consumerGroup consumer group err
	newConsumerGroupFunc = func(parent, path string, q FanOutQueue) (ConsumerGroup, error) {
		return nil, fmt.Errorf("err")
	}
	fq, err = NewFanOutQueue(dir, 1024)
	assert.Error(t, err)
	assert.Nil(t, fq)

	// case 6: init consumerGroup consumer group success
	newConsumerGroupFunc = func(parent, path string, q FanOutQueue) (ConsumerGroup, error) {
		return NewMockConsumerGroup(ctrl), nil
	}
	fq, err = NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	assert.NotNil(t, fq)
}

func TestFanOutQueue_GetOrCreateConsumerGroup(t *testing.T) {
	dir := filepath.Join(t.TempDir(), t.Name())

	defer func() {
		newConsumerGroupFunc = NewConsumerGroup
	}()

	fq, err := NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	assert.NotNil(t, fq)
	// case 1: create consumer group err
	newConsumerGroupFunc = func(parent, path string, q FanOutQueue) (ConsumerGroup, error) {
		return nil, fmt.Errorf("err")
	}
	fo, err := fq.GetOrCreateConsumerGroup("group-1")
	assert.Error(t, err)
	assert.Nil(t, fo)

	newConsumerGroupFunc = NewConsumerGroup

	// case 2: create consumer group success
	fo, err = fq.GetOrCreateConsumerGroup("group-1")
	assert.NoError(t, err)
	assert.NotNil(t, fo)

	fo2, err := fq.GetOrCreateConsumerGroup("group-1")
	assert.NoError(t, err)
	assert.Equal(t, fo, fo2)

	foNames := fq.ConsumerGroupNames()
	assert.Equal(t, "group-1", foNames[0])
	fq.Close()
}

func TestFanoutQueue_StopConsumerGroup(t *testing.T) {
	dir := filepath.Join(t.TempDir(), t.Name())
	fq, err := NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	assert.NotNil(t, fq)

	cgName := "group-1"

	fo, err := fq.GetOrCreateConsumerGroup(cgName)
	assert.NoError(t, err)
	assert.NotNil(t, fo)

	for i := 0; i < 10; i++ {
		_ = fq.Queue().Put([]byte(fmt.Sprintf("test-%d", i)))
	}

	seq := fo.Consume()
	assert.Equal(t, int64(0), seq)

	foNames := fq.ConsumerGroupNames()
	assert.Equal(t, cgName, foNames[0])

	// stop consumer group
	fq.StopConsumerGroup(cgName)
	assert.Empty(t, fq.ConsumerGroupNames())

	// reopen consumer group , can continue consume data
	fo, err = fq.GetOrCreateConsumerGroup(cgName)
	assert.NoError(t, err)
	assert.NotNil(t, fo)

	seq = fo.Consume()
	assert.Equal(t, int64(1), seq)
	foNames = fq.ConsumerGroupNames()
	assert.Equal(t, cgName, foNames[0])

	fq.Close()
}

func TestFanOutQueue_Sync(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := filepath.Join(t.TempDir(), t.Name())

	defer ctrl.Finish()

	fq, err := NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		err = fq.Queue().Put([]byte("12345"))
		assert.NoError(t, err)
	}

	assert.Equal(t, int64(-1), fq.Queue().AcknowledgedSeq())
	// case 1: sync with empty consume group
	fq.Sync()
	assert.Equal(t, int64(-1), fq.Queue().AcknowledgedSeq())
	fo1, err := fq.GetOrCreateConsumerGroup("group-1")
	assert.NoError(t, err)
	fo1.Consume()       //0
	fo1.Consume()       //1
	s1 := fo1.Consume() //2
	fo2, err := fq.GetOrCreateConsumerGroup("group-2")
	assert.NoError(t, err)

	s2 := fo2.Consume() //0

	// case 2: sync and ack min consume sequence, but consumer not ack
	fq.Sync()
	assert.Equal(t, int64(-1), fq.Queue().AcknowledgedSeq())
	// case 3: ack and sync
	fo1.Ack(s1)
	fo2.Ack(s2)
	fq.Sync()
	assert.Equal(t, int64(0), fq.Queue().AcknowledgedSeq())
	// case 3: sync err
	s2 = fo2.Consume() //1
	fo3 := fo2.(*consumerGroup)
	metaPage := page.NewMockMappedPage(ctrl)
	fo3.metaPage = metaPage
	metaPage.EXPECT().PutUint64(gomock.Any(), gomock.Any()).MaxTimes(2)
	metaPage.EXPECT().Sync().Return(fmt.Errorf("err"))
	fo2.Ack(s2)
	fq.Close()
}

func TestFanOutQueue_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	dir := filepath.Join(t.TempDir(), t.Name())

	defer ctrl.Finish()

	fq, err := NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	fo, err := fq.GetOrCreateConsumerGroup("f1")
	assert.NoError(t, err)
	assert.NotNil(t, fo)

	fo1 := fo.(*consumerGroup)
	pageFct := page.NewMockFactory(ctrl)
	metaPageFct := fo1.metaPageFct
	defer func() {
		_ = metaPageFct.Close()
	}()
	fo1.metaPageFct = pageFct
	pageFct.EXPECT().Close().Return(fmt.Errorf("err"))

	fq.Close()
}

func TestFanOutQueue_multiple_consumer(t *testing.T) {
	dir := filepath.Join(t.TempDir(), t.Name())

	fq, err := NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)
	// put data
	for i := 0; i < 100; i++ {
		msg := []byte(fmt.Sprintf("msg-%d", i))
		err = fq.Queue().Put(msg)
		assert.NoError(t, err)
	}

	// consumer group 1
	consumeMsg(t, fq, "f1", 100)
	// consumer group 2
	consumeMsg(t, fq, "f2", 100)

	fq.Close()
}

func TestFanOutQueue_SetAppendedSeq(t *testing.T) {
	dir := filepath.Join(t.TempDir(), t.Name())

	fq, err := NewFanOutQueue(dir, 1024)
	assert.NoError(t, err)

	// put data
	for i := 0; i < 100; i++ {
		msg := []byte(fmt.Sprintf("msg-%d", i))
		err = fq.Queue().Put(msg)
		assert.NoError(t, err)
	}

	// consumer group 1
	consumeMsg(t, fq, "f1", 100)
	f2, _ := fq.GetOrCreateConsumerGroup("f2")
	// set new append seq
	fq.SetAppendedSeq(200)
	// put msg
	msg := []byte(fmt.Sprintf("msg-%d", 200))
	err = fq.Queue().Put(msg)
	assert.NoError(t, err)

	f1, _ := fq.GetOrCreateConsumerGroup("f1")
	fseq := f1.Consume()
	assert.Equal(t, fseq, int64(201))
	fmsg, err := f1.Queue().Queue().Get(fseq)
	assert.NoError(t, err)
	assert.Equal(t, []byte(fmt.Sprintf("msg-%d", 200)), fmsg)

	fseq = f2.Consume()
	assert.Equal(t, fseq, int64(201))
	fmsg, err = f2.Queue().Queue().Get(fseq)
	assert.NoError(t, err)
	assert.Equal(t, []byte(fmt.Sprintf("msg-%d", 200)), fmsg)

	fq.Close()
}

func TestFanOutQueue_concurrent_read(t *testing.T) {
	dir := filepath.Join(t.TempDir(), t.Name())

	msgSize := 1024
	dataFileSize := int64(512)
	// random text
	bytesSli := make([][]byte, msgSize)

	for i := range bytesSli {
		bytesSli[i] = []byte(randomString(rand.Intn(10) + 1))
	}

	fq, err := NewFanOutQueue(dir, dataFileSize)
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

	assert.Equal(t, int64(msgSize-1), fq.Queue().AppendedSeq())

	// wait for background deleting
	time.Sleep(2 * time.Second)

	fq.Close()

	// reload
	fq2, err := NewFanOutQueue(dir, dataFileSize)
	assert.NoError(t, err)

	assert.Equal(t, int64(msgSize)-1, fq2.Queue().AppendedSeq())

	assert.Equal(t, len(fq2.ConsumerGroupNames()), readConcurrent)

	for i := 0; i < readConcurrent; i++ {
		fo, err := fq2.GetOrCreateConsumerGroup("fo-" + strconv.Itoa(i))
		assert.NoError(t, err)

		assert.Equal(t, int64(msgSize)-1, fo.ConsumedSeq())
		assert.Equal(t, int64(msgSize)-1, fo.AcknowledgedSeq())
	}
	fq2.Close()
}

func read(t *testing.T, raw [][]byte, fq FanOutQueue, name string) {
	fo, err := fq.GetOrCreateConsumerGroup(name)
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

		bys, err := fo.Queue().Queue().Get(seq)
		assert.NoError(t, err)
		assert.Equal(t, raw[int(seq)], bys)

		counter++

		fo.Ack(seq)
	}
}

func write(t *testing.T, raw [][]byte, fq FanOutQueue) {
	for _, bys := range raw {
		err := fq.Queue().Put(bys)
		assert.NoError(t, err)
	}
}

func consumeMsg(t *testing.T, fq FanOutQueue, consumerGroup string, msgCount int) {
	f1, err := fq.GetOrCreateConsumerGroup(consumerGroup)
	assert.NoError(t, err)

	for i := 0; i < msgCount; i++ {
		fseq := f1.Consume()
		assert.Equal(t, fseq, int64(i))
		fmsg, err := f1.Queue().Queue().Get(int64(i))
		assert.NoError(t, err)
		assert.Equal(t, []byte(fmt.Sprintf("msg-%d", i)), fmsg)
	}
}
