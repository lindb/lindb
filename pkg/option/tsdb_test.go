package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngineOption_Validation(t *testing.T) {
	engine := EngineOption{Interval: "ad"}
	assert.NotNil(t, engine.Validation())
	engine = EngineOption{Interval: "-10s"}
	assert.NotNil(t, engine.Validation())
	engine = EngineOption{Interval: "10s"}
	assert.Nil(t, engine.Validation())
	engine = EngineOption{Interval: "10s", Rollup: []string{"10s", "1m", "aa"}}
	assert.NotNil(t, engine.Validation())
	engine = EngineOption{Interval: "10s", Rollup: []string{"10s", "1m", "1h"}}
	assert.NotNil(t, engine.Validation())
	engine = EngineOption{Interval: "10s", Rollup: []string{"1s", "1m", "1h"}}
	assert.NotNil(t, engine.Validation())
	engine = EngineOption{Interval: "10s", Rollup: []string{"20s", "1m", "1h"}}
	assert.Nil(t, engine.Validation())
	engine = EngineOption{Interval: "10s", Rollup: []string{"10s", "1m", "1h"}, Ahead: "aa"}
	assert.NotNil(t, engine.Validation())
	engine = EngineOption{Interval: "10s", Rollup: []string{"10s", "1m", "1h"}, Behind: "aa"}
	assert.NotNil(t, engine.Validation())
	engine = EngineOption{Interval: "10s", Rollup: []string{"20s", "1m", "1h"}, Behind: "10h", Ahead: "1h"}
	assert.Nil(t, engine.Validation())
}
