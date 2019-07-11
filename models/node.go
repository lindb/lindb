package models

import (
	"fmt"
)

// Node represents the basic info of server
type Node struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

// String returns node info string
func (n *Node) String() string {
	return fmt.Sprintf("%s:%d", n.IP, n.Port)
}
