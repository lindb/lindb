package queue

import (
	"math/rand"
	"os"
	"path"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestOneFanOut(t *testing.T) {
	dir := path.Join(os.TempDir(), "fanOut")
	// remove dir to avoid influence of the previous run of test
	if err := os.RemoveAll(dir); err != nil {
		t.Error(err)
	}

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(fq.FanOutNames()), 0, "len(m)")

	f1, err := fq.GetOrCreateFanOut("f1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, f1.Name(), "f1")
	assert.Equal(t, f1.HeadSeq(), int64(0))
	assert.Equal(t, f1.TailSeq(), int64(0))
	assert.Equal(t, f1.Consume(), SeqNoNewMessageAvailable)
	assert.Equal(t, f1.Pending(), int64(0))

	assert.Equal(t, fq.HeadSeq(), int64(0))
	assert.Equal(t, fq.TailSeq(), int64(0))

	// msg 0
	msg := []byte("123")
	seq, err := fq.Append(msg)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, seq, int64(0))
	assert.Equal(t, f1.Pending(), int64(1))

	fseq := f1.Consume()
	assert.Equal(t, fseq, int64(0))
	assert.Equal(t, f1.HeadSeq(), int64(1))
	assert.Equal(t, f1.TailSeq(), int64(0))
	assert.Equal(t, f1.Pending(), int64(0))

	assert.Equal(t, f1.HeadSeq(), int64(1))
	assert.Equal(t, f1.TailSeq(), int64(0))

	fmsg, err := f1.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fmsg, msg)

	assert.Equal(t, f1.Consume(), SeqNoNewMessageAvailable)

	// msg1, msg2
	msg1 := []byte("456")
	msg2 := []byte("789")

	seq, err = fq.Append(msg1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seq, int64(1))

	assert.Equal(t, fq.HeadSeq(), int64(2))

	seq, err = fq.Append(msg2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seq, int64(2))
	assert.Equal(t, fq.HeadSeq(), int64(3))

	assert.Equal(t, f1.Pending(), int64(2))

	fseq = f1.Consume()
	assert.Equal(t, fseq, int64(1))
	assert.Equal(t, f1.HeadSeq(), int64(2))

	fmsg, err = f1.Get(fseq)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmsg, msg1)

	f1.Ack(fseq)
	assert.Equal(t, f1.TailSeq(), fseq)

	assert.Equal(t, fq.TailSeq(), fseq)

	fseq = f1.Consume()
	assert.Equal(t, fseq, int64(2))

	fmsg, err = f1.Get(fseq)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fmsg, msg2)
	f1.Ack(fseq)
	assert.Equal(t, f1.TailSeq(), fseq)
	assert.Equal(t, fq.TailSeq(), fseq)
	assert.Equal(t, f1.Pending(), int64(0))

	fq.Close()

}

func TestFanOut_SetHeadSeq(t *testing.T) {
	dir := path.Join(os.TempDir(), "fanOut")
	// remove dir to avoid influence of the previous run of test
	if err := os.RemoveAll(dir); err != nil {
		t.Error(err)
	}

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	f1, err := fq.GetOrCreateFanOut("f1")
	if err != nil {
		t.Fatal(err)
	}

	if err := f1.SetHeadSeq(1); err == nil {
		t.Fatal("should be error")
	}

	if _, err := fq.Append([]byte("123")); err != nil {
		t.Fatal(err)
	}

	if _, err := fq.Append([]byte("456")); err != nil {
		t.Fatal(err)
	}

	seq := f1.Consume()
	assert.Equal(t, seq, int64(0))

	seq = f1.Consume()
	assert.Equal(t, seq, int64(1))

	if err := f1.SetHeadSeq(0); err != nil {
		t.Fatal(err)
	}

	seq = f1.Consume()
	assert.Equal(t, seq, int64(0))

	seq = f1.Consume()
	assert.Equal(t, seq, int64(1))

	f1.Ack(1)

	if err := f1.SetHeadSeq(0); err == nil {
		t.Fatal("should be error")
	}
}

func TestMultipleFanOut(t *testing.T) {
	dir := path.Join(os.TempDir(), "fanOut")
	// remove dir to avoid influence of the previous run of test
	if err := os.RemoveAll(dir); err != nil {
		t.Error(err)
	}

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	f1, err := fq.GetOrCreateFanOut("f1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, f1.Consume(), SeqNoNewMessageAvailable)

	f2, err := fq.GetOrCreateFanOut("f2")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, f2.Consume(), SeqNoNewMessageAvailable)

	msg := []byte("123")

	seq, err := fq.Append(msg)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seq, int64(0))

	fseq := f1.Consume()
	assert.Equal(t, fseq, int64(0))

	fmsg, err := f1.Get(0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmsg, msg)

	fseq = f2.Consume()
	assert.Equal(t, fseq, int64(0))

	fmsg, err = f2.Get(0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmsg, msg)

	fq.Close()

}

func TestConcurrentRead(t *testing.T) {
	dir := path.Join(os.TempDir(), "fanout_concurrent")
	// remove dir to avoid influence of the previous run of test
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	msgSize := 1024
	dataFileSize := 512
	// random text
	bytesSli := make([][]byte, msgSize)

	for i := range bytesSli {
		bytesSli[i] = []byte(randomString(rand.Intn(10) + 1))
	}

	fq, err := NewFanOutQueue(dir, dataFileSize, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	readConcurrent := 10
	wg := &sync.WaitGroup{}

	for i := 0; i < readConcurrent; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()
			read(t, bytesSli, fq, "fo-"+strconv.Itoa(seq))
		}(i)
	}

	wg.Add(1)
	go func() {
		wg.Done()
		write(t, bytesSli, fq)
	}()

	wg.Wait()

	assert.Equal(t, fq.HeadSeq(), int64(msgSize))
	assert.Equal(t, fq.TailSeq(), int64(msgSize-1))

	// wait for background deleting
	time.Sleep(2 * time.Second)
	seg, err := fq.GetSegment(fq.TailSeq())
	if err != nil {
		t.Fatal(err)
	}

	_, err = fq.GetSegment(seg.Begin() - 1)
	if err == nil {
		t.Fatal(err, "should be deleted")
	}

	fq.Close()

	// reload
	fq2, err := NewFanOutQueue(dir, dataFileSize, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fq2.HeadSeq(), int64(msgSize))
	assert.Equal(t, fq2.TailSeq(), int64(msgSize-1))

	assert.Equal(t, len(fq2.FanOutNames()), readConcurrent)
	for i := 0; i < readConcurrent; i++ {
		fo, err := fq2.GetOrCreateFanOut("fo-" + strconv.Itoa(i))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fo.HeadSeq(), int64(msgSize))
		assert.Equal(t, fo.TailSeq(), int64(msgSize-1))
	}

	fq2.Close()

}

func read(t *testing.T, raw [][]byte, fq FanOutQueue, name string) {
	fo, err := fq.GetOrCreateFanOut(name)
	if err != nil {
		t.Error(err)
	}

	counter := fo.TailSeq()
	for {

		if counter == int64(len(raw)) {
			return
		}
		seq := fo.Consume()
		//fmt.Printf("fanout:%s, seq:%d\n", name, seq)
		if seq == SeqNoNewMessageAvailable {
			time.Sleep(time.Millisecond)
			continue
		}

		assert.Equal(t, seq, counter)

		bys, err := fo.Get(seq)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, bys, raw[int(seq)])

		counter++

		fo.Ack(seq)
	}
}

func write(t *testing.T, raw [][]byte, fq FanOutQueue) {
	for i, bys := range raw {
		seq, err := fq.Append(bys)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, int64(i), seq)
	}
}
