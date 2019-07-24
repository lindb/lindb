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

	"github.com/coreos/etcd/pkg/stringutil"
)

func TestOneFanOut(t *testing.T) {
	dir := path.Join(os.TempDir(), "fanOut")

	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}

	}()

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
	assert.Equal(t, f1.Consume(), int64(SeqNoNewMessageAvailable))

	assert.Equal(t, fq.HeadSeq(), int64(0))
	assert.Equal(t, fq.TailSeq(), int64(0))

	// msg 0
	msg := []byte("123")
	seq, err := fq.Append(msg)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, seq, int64(0))

	fseq := f1.Consume()
	assert.Equal(t, fseq, int64(0))
	assert.Equal(t, f1.HeadSeq(), int64(1))
	assert.Equal(t, f1.TailSeq(), int64(0))

	assert.Equal(t, f1.HeadSeq(), int64(1))
	assert.Equal(t, f1.TailSeq(), int64(0))

	fmsg, err := f1.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fmsg, msg)

	assert.Equal(t, f1.Consume(), int64(SeqNoNewMessageAvailable))

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

	fq.Close()

}

func TestMultipleFanOut(t *testing.T) {
	dir := path.Join(os.TempDir(), "fanOut")

	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}

	}()

	fq, err := NewFanOutQueue(dir, 1024, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	f1, err := fq.GetOrCreateFanOut("f1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, f1.Consume(), int64(SeqNoNewMessageAvailable))

	f2, err := fq.GetOrCreateFanOut("f2")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, f2.Consume(), int64(SeqNoNewMessageAvailable))

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

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()

	msgSize := 1024
	dataFileSize := 512
	// random text
	bytesSli := make([][]byte, msgSize, msgSize)

	for i := range bytesSli {
		bytesSli[i] = []byte(stringutil.RandomStrings(uint(rand.Intn(10)+1), 1)[0])
	}

	fq, err := NewFanOutQueue(dir, dataFileSize, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	conc := 10
	wg := &sync.WaitGroup{}
	wg.Add(conc)
	wg.Add(1)

	for i := 0; i < conc; i++ {
		go read(t, bytesSli, fq, "fo-"+strconv.Itoa(i), wg)
	}

	go write(t, bytesSli, fq, wg)

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

	assert.Equal(t, len(fq2.FanOutNames()), conc)
	for i := 0; i < conc; i++ {
		//go read(t, bytesSli, fq, "fo-"+strconv.Itoa(i), wg)
		fo, err := fq2.GetOrCreateFanOut("fo-" + strconv.Itoa(i))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fo.HeadSeq(), int64(msgSize))
		assert.Equal(t, fo.TailSeq(), int64(msgSize-1))
	}

	fq2.Close()

}

func read(t *testing.T, raw [][]byte, fq FanOutQueue, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fo, err := fq.GetOrCreateFanOut(name)
	if err != nil {
		t.Error(err)
	}

	counter := fo.TailSeq()
	//t.Logf("fanout %s consume from %d", fo.Name(), counter)
	for {

		if counter == int64(len(raw)) {
			return
		}
		seq := fo.Consume()
		if seq == -1 {
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

func write(t *testing.T, raw [][]byte, fq FanOutQueue, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, bys := range raw {
		seq, err := fq.Append(bys)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, int64(i), seq)
	}
}
