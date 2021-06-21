// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package models

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// NodeID represents node identifier.
type NodeID int

// Node represents the basic info of server
type Node struct {
	IP       string `json:"ip"`
	Port     uint16 `json:"port"`
	HTTPPort uint16 `json:"httpPort"`
	HostName string `json:"hostName"`
}

// Indicator returns return node indicator's string
func (n *Node) Indicator() string {
	return fmt.Sprintf("%s:%d", n.IP, n.Port)
}

// ParseNode parses Node from indicator,
// if indicator is not in the form [ip]:port  or port is not valid num, return error.
func ParseNode(indicator string) (*Node, error) {
	index := strings.Index(indicator, ":")
	if index < 0 {
		return nil, fmt.Errorf("indicator(%s) is not in the format [ip]:port", indicator)
	}

	ipStr := indicator[:index]
	if ip := net.ParseIP(ipStr); ip == nil {
		return nil, fmt.Errorf("indicator(%s) contains a invalid ip address", indicator)
	}

	portStr := indicator[index+1:]
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, err
	}
	return &Node{
		IP:   indicator[:index],
		Port: uint16(port),
	}, nil
}

// Master represents master basic info
type Master struct {
	Node      Node  `json:"node"`
	ElectTime int64 `json:"electTime"`
}

// ActiveNode represents active node include online time
type ActiveNode struct {
	Version    string `json:"version"`
	Node       Node   `json:"node"`
	OnlineTime int64  `json:"onlineTime"` // node online time(millisecond)
}
