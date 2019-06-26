package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostIP(t *testing.T) {
	ip, err := GetHostIP()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip)
	fmt.Println("ip:" + ip)
}
