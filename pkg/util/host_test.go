package util

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostIP(t *testing.T) {
	ip, err := GetHostIP()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip)
	fmt.Println("ip:" + ip)
}

func TestGetHostName(t *testing.T) {
	fmt.Println(GetHostName())
	assert.NotNil(t, GetHostName())
	defer func() { osHostname = os.Hostname }()
	osHostname = func() (string, error) { return "", errors.New("fail") }
	assert.Equal(t, "unknown", GetHostName())
}
