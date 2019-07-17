package rpc

import (
	"time"

	"google.golang.org/grpc"
)

type ClientConfig struct {
	Address string
	Timeout time.Duration
}

type Client interface {
	Conn() *grpc.ClientConn
	Close() error
}

type client struct {
	conn *grpc.ClientConn
	cfg  ClientConfig
}

func NewClient(cfg ClientConfig) (Client, error) {
	c := &client{
		cfg: cfg,
	}
	conn, err := grpc.Dial(cfg.Address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return c, nil
}

func (c *client) Conn() *grpc.ClientConn {
	return c.conn
}
func (c *client) Close() error {
	return nil
}
