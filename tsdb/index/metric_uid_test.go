package index

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
)

var testKVPath = "../test"

func Test_MeasurementAdd(t *testing.T) {
	//TODO need modify test case
	fileutil.RemoveDir(testKVPath)
	defer fileutil.RemoveDir(testKVPath)
	option := kv.DefaultStoreOption(testKVPath)
	var indexStore, _ = kv.NewStore("index", option)
	family, _ := indexStore.CreateFamily("measurement", kv.FamilyOption{})

	measurementUID := NewMetricUID(family)
	for i := 0; i < 100; i++ {
		measurementUID.GetOrCreateMetricID(fmt.Sprintf("%s%d", "key-", i), true)
	}
	measurementUID.Flush()

	for i := 98; i < 200; i++ {
		measurementUID.GetOrCreateMetricID(fmt.Sprintf("%s%d", "key-", i), true)
	}

	for i := 0; i < 100; i++ {
		id, ok := measurementUID.GetOrCreateMetricID(fmt.Sprintf("%s%d", "key-", i), false)
		if ok {
			assert.Equal(t, uint32(i+1), id)
		}
	}

	_ = indexStore.Close()
}
