package rpc

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/eleme/lindb/rpc/proto/broker"
	"github.com/eleme/lindb/rpc/proto/common"
)

type BrokerClient interface {
	Init() error
	WritePoints(request *common.Request) (*common.Response, error)
	Close() error
}

type brokerClient struct {
	conn    *grpc.ClientConn
	client  broker.BrokerServiceClient
	address string
	timeout time.Duration
}

func NewBrokerClient(address string, timeout time.Duration) BrokerClient {
	return &brokerClient{
		address: address,
		timeout: timeout,
	}
}

func (bc *brokerClient) Init() error {
	conn, err := grpc.Dial(bc.address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	bc.conn = conn

	bc.client = broker.NewBrokerServiceClient(conn)

	return nil
}

func (bc *brokerClient) WritePoints(request *common.Request) (*common.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bc.timeout)
	defer cancel()
	return bc.client.WritePoints(ctx, request)
}

func (bc *brokerClient) Close() error {
	if bc.conn != nil {
		return bc.conn.Close()
	}
	return nil
}
