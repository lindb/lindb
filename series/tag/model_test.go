package tag

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tags(t *testing.T) {
	var tags = Tags{}
	assert.Len(t, tags.AppendHashKey(nil), 0)
	tags = append(tags, NewTag([]byte("ip"), []byte("1.1.1.1")),
		NewTag([]byte("zone"), []byte("sh")),
		NewTag([]byte("host"), []byte("test")))
	assert.Equal(t, 23, tags.Size())
	assert.False(t, tags.needsEscape())
	assert.Equal(t, ",ip=1.1.1.1,zone=sh,host=test", string(tags.AppendHashKey(nil)))

	tags = append(tags, NewTag([]byte("x x"), []byte("y,y")))
	sort.Sort(tags)
	assert.True(t, tags.needsEscape())
	assert.Equal(t, ",host=test,ip=1.1.1.1,x\\ x=y\\,y,zone=sh", string(tags.AppendHashKey(nil)))
}
