package broker

import (
	"context"
	"time"

	"github.com/eleme/lindb/rpc"
	"github.com/eleme/lindb/rpc/pkg/batch"
	brokerpb "github.com/eleme/lindb/rpc/pkg/broker"
	"github.com/eleme/lindb/rpc/pkg/common"
)

/***
* procedures to add RPC service, take add WritPoints to BrokerService as an example
* * add request type WritePoints in rpc/pb/common.proto
* * define WritePointsRequest and WritePointsResponse, writePoints rpc method in rpc/pb/broker.proto,
* * add WritePointsRequest, WritePointsResponse in rpc/pb/batch.proto BatchRequest.Request BatchResponse.Response
* * run `make pb` to regenerate pb files in rpc/pkg/broker
* * define broker_client.go wrappers BaseClient in rpc/broker.go,
* * RegisterToBatchRequestFunc RegisterFromBatchFunc in func init()
* * send stream requests through the baseClient, send direct request through grpc BrokerServiceClient
* * define broker_server.go implements rpc.Server and grpc service interface
 */

func init() {
	rpc.RegisterToBatchRequestFunc(common.RequestType_WritePoints, func(req *rpc.Request) *batch.BatchRequest_Request {
		return &batch.BatchRequest_Request{
			RequestTyp: common.RequestType_WritePoints,
			Request: &batch.BatchRequest_Request_WritePoints{
				WritePoints: req.WritePointRequest,
			},
		}
	})

	rpc.RegisterFromBatchFunc(common.RequestType_WritePoints, func(res *batch.BatchResponse_Response) *rpc.Response {
		return &rpc.Response{
			Type:                common.RequestType_WritePoints,
			WritePointsResponse: res.GetWritePoints(),
		}
	})
}

type Client interface {
	WritePoints(points []*brokerpb.Point) error
	WritePointsStream(points []*brokerpb.Point) error
	Close()
}

type brokerClient struct {
	baseCli rpc.BaseClient
	address string
	timeout time.Duration
}

func NewBrokerClient(address string, timeout time.Duration) Client {
	cli := rpc.NewBaseClient()
	return &brokerClient{
		baseCli: cli,
		address: address,
		timeout: timeout,
	}
}

func (bc *brokerClient) WritePoints(points []*brokerpb.Point) error {
	conn, err := bc.baseCli.GetConn(bc.address)
	if err != nil {
		return err
	}
	writePointRequest :=
		&brokerpb.WritePointsRequest{
			Context: rpc.BuildRequestContext(common.RequestType_WritePoints),
			Points:  points,
		}

	// todo cache BrokerServiceClient
	cli := brokerpb.NewBrokerServiceClient(conn)
	_, err = cli.WritePoints(context.TODO(), writePointRequest)
	return err
}

func (bc *brokerClient) WritePointsStream(points []*brokerpb.Point) error {
	req := &rpc.Request{
		Type: common.RequestType_WritePoints,
		WritePointRequest: &brokerpb.WritePointsRequest{
			Context: rpc.BuildRequestContext(common.RequestType_WritePoints),
			Points:  points,
		},
	}

	_, err := bc.baseCli.SendStreamRequest(context.TODO(), bc.address, req, bc.timeout)

	return err
}

func (bc *brokerClient) Close() {
	bc.baseCli.Close()
}
