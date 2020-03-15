package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUser struct {
	Name string
}

func TestJSONCodec(t *testing.T) {
	user := mockUser{Name: "LinDB"}
	data := JSONMarshal(&user)
	user1 := mockUser{}
	err := JSONUnmarshal(data, &user1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, user, user)
	err = JSONUnmarshal([]byte{1, 1, 1}, &user1)
	assert.NotNil(t, err)
}

func Test_JSONMarshal(t *testing.T) {
	assert.Len(t, JSONMarshal(make(chan struct{}, 1)), 0)
	assert.True(t, len(JSONMarshal(nil)) > 0)
}
