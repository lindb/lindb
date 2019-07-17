package rpc

type ClientPool interface {
	Choose(address string) *Client
	Add(cfg ClientConfig) error
	Remove(address string) error
}
