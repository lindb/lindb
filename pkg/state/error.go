package state

import (
	"fmt"

	etcdcliv3 "github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
)

var (
	// ErrWatchFailed indicates the watch failed.
	ErrWatchFailed = fmt.Errorf("etcd watch returns a nil chan")
	// ErrNoKey indicates the key does not exist.
	ErrNoKey = fmt.Errorf("etcd has no such key")
	// ErrTxnFailed indicates the txn failed.
	ErrTxnFailed = fmt.Errorf("role changed or target revision mismatch")
	// ErrTxnConvert transaction covert failed.
	ErrTxnConvert = fmt.Errorf("cannot covert etcd transaction")
)

// TxnErr converts txn response and error into one error.
func TxnErr(resp *etcdcliv3.TxnResponse, err error) error {
	if err != nil {
		return errors.WithStack(err)
	}
	if !resp.Succeeded {
		return ErrTxnFailed
	}
	return nil
}
