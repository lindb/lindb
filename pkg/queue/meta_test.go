package queue

import (
	"os"
	"path"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestMeta(t *testing.T) {
	tmpFilePath := path.Join(os.TempDir(), "testMeta.meta")
	defer func() {
		if err := os.Remove(tmpFilePath); err != nil {
			t.Error(err)
		}
	}()

	meta, err := NewMeta(tmpFilePath, 16)
	if err != nil {
		t.Fatal(err)
	}

	meta.WriteInt64(0, 123)
	meta.WriteInt64(8, 456)

	if err := meta.Sync(); err != nil {
		t.Fatal(err)
	}

	if err := meta.Close(); err != nil {
		t.Fatal(err)
	}

	meta, err = NewMeta(tmpFilePath, 16)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, meta.ReadInt64(0), int64(123))
	assert.Equal(t, meta.ReadInt64(8), int64(456))
}

func TestNewMetaError(t *testing.T) {
	tmpFilePath := path.Join(os.TempDir(), "testMeta.meta")
	defer func() {
		if err := os.Remove(tmpFilePath); err != nil {
			t.Error(err)
		}
	}()

	_, err := NewMeta(tmpFilePath, -1)
	if err == nil {
		t.Fatal("should be error")
	}

}
