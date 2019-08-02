package index

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/util"
)

func TestFieldUid_GetOrCreateFieldId(t *testing.T) {
	defer fileutil.RemoveDir("../test")
	fieldUID := NewFieldUID(initFieldFamily())

	for i := 1; i < 10; i++ {
		for j := 1; j < 100; j++ {
			id, _ := fieldUID.GetOrCreateFieldID(uint32(i), fmt.Sprintf("%s%d", "field-", j), field.Type(i))
			high, low := util.IntToShort(id)
			assert.Equal(t, i, int(high))
			assert.Equal(t, j, int(low))
		}
	}
	err := fieldUID.Flush()
	assert.Equal(t, nil, err)

	for i := 1; i < 10; i++ {
		for j := 1; j < 100; j++ {
			id := fieldUID.GetFieldID(uint32(i), fmt.Sprintf("%s%d", "field-", j))
			if id == NotFoundFieldID {
				fmt.Println(i, j)
			}
			high, low := util.IntToShort(id)
			assert.Equal(t, i, int(high))
			assert.Equal(t, j, int(low))
		}
	}
	fmt.Println("success")
}

func initFieldFamily() kv.Family {
	option := kv.DefaultStoreOption("../test")

	var indexStore, _ = kv.NewStore("index", option)

	family, _ := indexStore.CreateFamily("field", kv.FamilyOption{})
	return family
}
