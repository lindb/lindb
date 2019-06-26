package util

import (
	"errors"
	"net"
	"sync"
)

var (
	once sync.Once
	host hostInfo
)

// hostInfo defines host basic info, if cannot get host info returns error
type hostInfo struct {
	hostIP string
	err    error
}

// extractHostInfo extracts host info, just do it once
func extractHostInfo() {
	once.Do(func() {
		ifaces, err := net.Interfaces()
		if err != nil {
			host.err = err
			return
		}
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				// interface is down or loopback
				continue
			}
			addrs, err := iface.Addrs()
			if err != nil {
				host.err = err
				return
			}
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip == nil || ip.IsLoopback() {
					continue
				}
				ip = ip.To4()
				if ip == nil {
					// not an ipv4 address
					continue
				}
				host.hostIP = ip.String()
				return
			}
		}
		host.err = errors.New("cannot extract host info")
	})
}

// GetHostIP returns current host ip address
func GetHostIP() (string, error) {
	extractHostInfo()
	return host.hostIP, host.err
}
