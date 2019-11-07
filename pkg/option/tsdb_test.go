package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DatabaseOption_Validate(t *testing.T) {
	databaseOption := DatabaseOption{Interval: "ad"}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "-10s"}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s"}
	assert.Nil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"10s", "1m", "aa"}}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"10s", "1m", "1h"}}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"1s", "1m", "1h"}}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"20s", "1m", "1h"}}
	assert.Nil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"10s", "1m", "1h"}, Ahead: "aa"}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"10s", "1m", "1h"}, Behind: "aa"}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"20s", "1m", "1h"}, Behind: "10h", Ahead: "1h"}
	assert.Nil(t, databaseOption.Validate())
}
