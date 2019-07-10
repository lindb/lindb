package signature

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_Aes(t *testing.T) {
	cid, _ := BuildNewUserInfo("etrace", "etrace")
	key, _ := AesDecrypt(cid)
	assert.Equal(t, "etrace/etrace", key)
}
