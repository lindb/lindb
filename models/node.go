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

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/lindb/lindb/pkg/timeutil"
)

// NodeID represents node identifier.
type NodeID int

func (id NodeID) Int() int       { return int(id) }
func (id NodeID) String() string { return strconv.Itoa(int(id)) }
func ParseNodeID(node string) NodeID {
	id, _ := strconv.Atoi(node)
	return NodeID(id)
}

// Node represents the node info in cluster(broker/storage).
type Node interface {
	// Indicator returns node indicator's string.
	Indicator() string
	HTTPAddress() string
}

type StatefulNode struct {
	StatelessNode

	ID NodeID `json:"id"`
}

// StatelessNodes represents stateless node list.
type StatelessNodes []StatelessNode

// ToTable returns stateless node list as table if it has value, else return empty string.
func (n StatelessNodes) ToTable() (int, string) {
	if len(n) == 0 {
		return 0, ""
	}
	writer := NewTableFormatter()
	writer.AppendHeader(table.Row{"Online time", "Host IP", "Host Name", "Port(HTTP/GRPC)", "Version"})
	for i := range n {
		r := n[i]
		writer.AppendRow(table.Row{
			timeutil.FormatTimestamp(r.OnlineTime, timeutil.DataTimeFormat2),
			r.HostIP, r.HostName, fmt.Sprintf("%d/%d", r.HTTPPort, r.GRPCPort), r.Version})
	}
	return len(n), writer.Render()
}

// StatelessNode represents stateless node basic info.
type StatelessNode struct {
	HostIP   string `json:"hostIp"`
	HostName string `json:"hostName"`
	GRPCPort uint16 `json:"grpcPort"`
	HTTPPort uint16 `json:"httpPort"`

	Version    string `json:"version"`
	OnlineTime int64  `json:"onlineTime"` // node online time(millisecond)
}

// Indicator returns node indicator's string.
func (n *StatelessNode) Indicator() string {
	return fmt.Sprintf("%s:%d", n.HostIP, n.GRPCPort)
}

func (n *StatelessNode) HTTPAddress() string {
	return fmt.Sprintf("http://%s:%d", n.HostIP, n.HTTPPort)
}

// ParseNode parses Node from indicator,
// if indicator is not in the form [ip]:port  or port is not valid num, return error.
func ParseNode(indicator string) (Node, error) {
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
	//TODO change base node info???
	return &StatelessNode{
		HostIP:   indicator[:index],
		GRPCPort: uint16(port),
	}, nil
}

// Master represents master basic info.
type Master struct {
	Node      *StatelessNode `json:"node"`
	ElectTime int64          `json:"electTime"`
}

// ToTable returns master info as table.
func (m *Master) ToTable() (int, string) {
	writer := NewTableFormatter()
	writer.AppendHeader(table.Row{"Desc", "Value"})
	writer.AppendRow(table.Row{"Elect Time", timeutil.FormatTimestamp(m.ElectTime, timeutil.DataTimeFormat2)})
	writer.AppendRow(table.Row{"Online Time", timeutil.FormatTimestamp(m.Node.OnlineTime, timeutil.DataTimeFormat2)})
	writer.AppendRow(table.Row{"Host IP", m.Node.HostIP})
	writer.AppendRow(table.Row{"Host Name", m.Node.HostName})
	writer.AppendRow(table.Row{"HTTP Port", m.Node.HTTPPort})
	writer.AppendRow(table.Row{"GRPC Port", m.Node.GRPCPort})
	return 1, writer.Render()
}
