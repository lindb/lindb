package ltoml

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Duration(t *testing.T) {
	assert.Equal(t, Duration(time.Minute).Duration(), time.Minute)

	marshalF := func(duration time.Duration) string {
		txt, _ := Duration(duration).MarshalText()
		return string(txt)
	}
	unmarshalF := func(txt string) time.Duration {
		var d Duration
		_ = d.UnmarshalText([]byte(txt))
		return d.Duration()
	}
	assert.Equal(t, "1m0s", marshalF(time.Minute))
	assert.Equal(t, "10s", marshalF(time.Second*10))

	assert.Equal(t, time.Second, unmarshalF("1s"))
	assert.Equal(t, time.Minute, unmarshalF("1m"))
	assert.Equal(t, time.Hour, unmarshalF("3600s"))

	assert.Zero(t, unmarshalF(""))
	assert.Zero(t, unmarshalF("1fs"))
}
