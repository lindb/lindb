package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/bufioutil"
)

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	bitWriter := bit.NewWriter(&buf)
	e := NewXOREncoder(bitWriter)
	e.Write(uint64(76))
	e.Write(uint64(50))
	e.Write(uint64(50))
	e.Write(uint64(999999999))
	e.Write(uint64(100))

	err := bitWriter.Flush()
	if err != nil {
		t.Fatal(err)
	}
	data := buf.Bytes()

	reader := bit.NewReader(bufioutil.NewBuffer(data))
	d := NewXORDecoder(reader)
	exceptIntValue(d, t, uint64(76))
	exceptIntValue(d, t, uint64(50))
	exceptIntValue(d, t, uint64(50))
	exceptIntValue(d, t, uint64(999999999))
	exceptIntValue(d, t, uint64(100))
}

func exceptIntValue(d *XORDecoder, t *testing.T, except uint64) {
	assert.True(t, d.Next())
	assert.Equal(t, except, d.Value())
}
