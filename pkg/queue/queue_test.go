package queue

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestOneSegment(t *testing.T) {
	dir := path.Join(os.TempDir(), "queue")

	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}

	}()

	q, err := NewQueue(dir, 1024, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, q.Size(), int64(0))
	assert.Equal(t, q.HeadSeq(), int64(0))
	assert.Equal(t, q.TailSeq(), int64(0))

	seq, err := q.Append([]byte("123"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, seq, int64(0))
	assert.Equal(t, q.Size(), int64(1))
	assert.Equal(t, q.HeadSeq(), int64(1))
	assert.Equal(t, q.TailSeq(), int64(0))

	seq, err = q.Append([]byte("456"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, seq, int64(1))
	assert.Equal(t, q.Size(), int64(2))
	assert.Equal(t, q.HeadSeq(), int64(2))
	assert.Equal(t, q.TailSeq(), int64(0))

	q.Ack(1)
	assert.Equal(t, q.Size(), int64(1))
	assert.Equal(t, q.TailSeq(), int64(1))

	q.Close()

}

func TestMultipleSegments(t *testing.T) {
	dir := path.Join(os.TempDir(), "queue")

	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}

	}()

	// interval 1 second for test
	q, err := NewQueue(dir, 10, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, q.Size(), int64(0))
	assert.Equal(t, q.HeadSeq(), int64(0))
	assert.Equal(t, q.TailSeq(), int64(0))

	// 1 segment
	seq, err := q.Append([]byte("0123456789"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, seq, int64(0))
	assert.Equal(t, q.Size(), int64(1))
	assert.Equal(t, q.HeadSeq(), int64(1))
	assert.Equal(t, q.TailSeq(), int64(0))

	q.Ack(0)
	assert.Equal(t, q.TailSeq(), int64(0))
	// wait enough time
	time.Sleep(2 * time.Second)
	if _, err := q.GetSegment(0); err != nil {
		t.Fatal(err, "segment[0,1) should exists")
	}

	// 2 segment
	seq, err = q.Append([]byte("a"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, seq, int64(1))
	assert.Equal(t, q.Size(), int64(2))
	assert.Equal(t, q.HeadSeq(), int64(2))
	assert.Equal(t, q.TailSeq(), int64(0))

	seq, err = q.Append([]byte("bcd"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, seq, int64(2))
	assert.Equal(t, q.Size(), int64(3))
	assert.Equal(t, q.HeadSeq(), int64(3))
	assert.Equal(t, q.TailSeq(), int64(0))

	q.Ack(1)
	assert.Equal(t, q.TailSeq(), int64(1))
	time.Sleep(2 * time.Second)
	if _, err := q.GetSegment(0); err == nil {
		t.Fatal(err, "segment[0,1) should not exists")
	}

	q.Close()
}
