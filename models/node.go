package models

import (
	"fmt"
)

// Node represents the basic info of server
type Node struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

// Indicator returns return node indicator's string
func (n *Node) Indicator() string {
	return fmt.Sprintf("%s:%d", n.IP, n.Port)
}

// Master represents master basic info
type Master struct {
	Node      Node  `json:"node"`
	ElectTime int64 `json:"electTime"`
}

// ActiveNode represents active node include online time
type ActiveNode struct {
	Node       Node  `json:"node"`
	OnlineTime int64 `json:"onlineTime"` // node online time(millisecond)
}
