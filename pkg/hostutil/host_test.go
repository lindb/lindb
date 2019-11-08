package hostutil

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostIP(t *testing.T) {
	ip, err := GetHostIP()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip)
	fmt.Println("ip:" + ip)
}

func Test_getHostInfo(t *testing.T) {
	defer func() {
		netInterfaces = net.Interfaces
	}()

	// mock err
	netInterfaces = func() (interfaces []net.Interface, e error) {
		return nil, fmt.Errorf("err")
	}
	host := getHostInfo()
	assert.Empty(t, host.hostIP)
	assert.Error(t, host.err)

	// mock empty
	netInterfaces = func() (interfaces []net.Interface, e error) {
		return nil, nil
	}
	host = getHostInfo()
	assert.Empty(t, host.hostIP)
	assert.Error(t, host.err)

	netInterfaces = func() (interfaces []net.Interface, e error) {
		return []net.Interface{
			{
				Name:  "mock_test",
				Flags: 1,
			},
		}, nil
	}
	_ = getHostInfo()
}
