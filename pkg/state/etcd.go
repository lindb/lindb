package state

import (
	"context"
	"fmt"

	etcd "go.etcd.io/etcd/clientv3"
)

// etcdRepository is repository based on etc storage
type etcdRepository struct {
	client *etcd.Client
}

// newETCDRepository creates a new repository based on etcd storage
func newETCDRepository(config interface{}) (Repository, error) {
	v, ok := config.(etcd.Config)
	if !ok {
		return nil, fmt.Errorf("config type is not etc.confit")
	}
	cli, err := etcd.New(v)
	if err != nil {
		return nil, fmt.Errorf("create etc client error:%s", err)
	}
	return &etcdRepository{
		client: cli,
	}, nil
}

// Get retrieves value for given key from etcd
func (r *etcdRepository) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get value failure for key[%s], error:%s", key, err)
	}
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("key[%s] not exist", key)
	}

	firstkv := resp.Kvs[0]
	if len(firstkv.Value) == 0 {
		return nil, fmt.Errorf("key[%s]'s value is empty", key)
	}
	return firstkv.Value, err
}

// Put puts a key-value pair into etcd
func (r *etcdRepository) Put(ctx context.Context, key string, val []byte) error {
	_, err := r.client.Put(ctx, key, string(val))
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes value for given key from etcd
func (r *etcdRepository) Delete(ctx context.Context, key string) error {
	_, err := r.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

// Close closes etcd client
func (r *etcdRepository) Close() error {
	return r.client.Close()
}
