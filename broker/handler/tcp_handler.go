package handler

import (
	"bufio"
	"net"

	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/field"
)

const (
	int32BytesLen = 4
)

type tcpHandler struct {
	channelManager replication.ChannelManager
}

func NewTCPHandler(cm replication.ChannelManager) rpc.TCPHandler {
	return &tcpHandler{channelManager: cm}
}

/**
tcp packet stream
packet size int32 4 bytes
bytes       []byte packet size bytes
*/
// Handles incoming requests.
func (h *tcpHandler) Handle(conn net.Conn) error {
	scanner := bufio.NewScanner(conn)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return
		}

		// try to extract complete data
		dataLen := len(data)
		if dataLen <= int32BytesLen {
			return
		}
		packageLen := int(stream.NewReader(data).ReadInt32())
		if dataLen < int32BytesLen+packageLen {
			return
		}
		return int32BytesLen + packageLen, data[int32BytesLen : int32BytesLen+packageLen], nil
	})

	for scanner.Scan() {
		data := scanner.Bytes()
		//// handler data
		//if len(data) == 0 {
		//	continue
		//}

		var metricList field.MetricList
		if err := metricList.Unmarshal(data); err != nil {
			return err
		}

		if err := h.channelManager.Write(&metricList); err != nil {
			return err
		}

	}
	return scanner.Err()
}
