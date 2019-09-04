package handler

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/field"
)

func TestTcpHandler_Handle(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cm := replication.NewMockChannelManager(ctl)

	h := NewTCPHandler(cm)
	tcpServer := rpc.NewTCPServer(":9000", h)

	go func() {
		if err := tcpServer.Start(); err != nil {
			fmt.Printf("DEBUG IGNORE %v", err)
		}
	}()

	time.Sleep(20 * time.Millisecond)
	wg := &sync.WaitGroup{}

	wg.Add(5)
	go func() {
		testMetrics(wg, t, cm, 1, nil)
	}()
	go func() {
		testMetrics(wg, t, cm, 2, errors.New("mock err"))
	}()
	go func() {
		testMetrics2(wg, t, cm, 2, nil)
	}()
	go func() {
		testWrongBytes1(wg, t, cm)
	}()
	go func() {
		testWrongBytes2(wg, t, cm)
	}()

	wg.Wait()
	// wait for server handler go routine
	time.Sleep(20 * time.Millisecond)
	tcpServer.Stop()

}

func testMetrics(wg *sync.WaitGroup, t *testing.T, cm *replication.MockChannelManager, value float64, mockErr error) {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		t.Fatal(err)
	}
	writeMetricList(wg, t, cm, conn, value, mockErr)

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}

func testMetrics2(wg *sync.WaitGroup, t *testing.T, cm *replication.MockChannelManager, value float64, mockErr error) {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)
	writeMetricList(wg, t, cm, conn, value, nil)
	writeMetricList(wg, t, cm, conn, value, mockErr)

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}

func writeMetricList(wg *sync.WaitGroup, t *testing.T, cm *replication.MockChannelManager, conn net.Conn, value float64, mockErr error) {
	defer wg.Done()
	metricList := buildMetricList(value)
	metricListBytes, err := metricList.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	writer := stream.NewBufferWriter(nil)
	writer.PutInt32(int32(len(metricListBytes)))
	writer.PutBytes(metricListBytes)

	bytes, err := writer.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	cm.EXPECT().Write(metricList).Return(mockErr)

	n, err := conn.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(bytes) {
		t.Fatal("should not happen")
	}
}

func testWrongBytes1(wg *sync.WaitGroup, t *testing.T, _ *replication.MockChannelManager) {
	defer wg.Done()
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		t.Fatal(err)
	}

	writer := stream.NewBufferWriter(nil)
	writer.PutByte(1)

	bytes, err := writer.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	n, err := conn.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(bytes) {
		t.Fatal("should not happen")
	}

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}

func testWrongBytes2(wg *sync.WaitGroup, t *testing.T, _ *replication.MockChannelManager) {
	defer wg.Done()
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		t.Fatal(err)
	}

	writer := stream.NewBufferWriter(nil)
	writer.PutInt32(2)
	writer.PutByte(1)

	bytes, err := writer.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	n, err := conn.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(bytes) {
		t.Fatal("should not happen")
	}

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}

func buildMetricList(value float64) *field.MetricList {
	return &field.MetricList{Database: "dal",
		Metrics: []*field.Metric{{
			Name:      "name",
			Timestamp: time.Now().Unix() * 1000,
			Tags:      map[string]string{"tagKey": "tagVal"},
			Fields: []*field.Field{{
				Name: "sum",
				Field: &field.Field_Sum{
					Sum: &field.Sum{
						Value: value,
					},
				},
			}},
		}}}
}

//// 发送数据用
//func TestWriteToBroker(t *testing.T) {
//	conn, err := net.Dial("tcp", ":9002")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	for i := 0; i < 1000; i++ {
//		ml := buildMetricList(float64(i))
//
//		metricListBytes, err := ml.Marshal()
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		writer := stream.NewBufferWriter(nil)
//		writer.PutInt32(int32(len(metricListBytes)))
//		writer.PutBytes(metricListBytes)
//
//		bytes, err := writer.Bytes()
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		if _, err := conn.Write(bytes); err != nil {
//			t.Fatal(err)
//		}
//	}
//
//	if err := conn.Close(); err != nil {
//		t.Fatal(err)
//	}
//}
