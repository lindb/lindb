package rest

import (
	"net/http"
	"testing"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"

	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/state"
)

func TestGetDatabase(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)
	clientConfig := etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	}
	_ = state.New("etcd", clientConfig)

	//create success
	originDatabase := option.Database{Name: "test", NumOfShard: 1, ReplicaFactor: 1}
	originHandlerTest := &httpHandlerTest{
		http.MethodPost,
		"/database",
		originDatabase,
		CreateOrUpdateDatabase,
		204,
		nil}
	testHTTPHandler(t, originHandlerTest)
	// get success
	successHandlerTest := &httpHandlerTest{
		http.MethodGet,
		"/database?databaseName=test",
		nil,
		GetDatabase,
		200,
		originDatabase}
	testHTTPHandler(t, successHandlerTest)
	// no database name
	noDatabaseNameHandlerTest := &httpHandlerTest{
		http.MethodGet,
		"/database/databaseName",
		nil,
		GetDatabase,
		500,
		nil}
	testHTTPHandler(t, noDatabaseNameHandlerTest)
	// wrong database name
	errorNameHandlerTest := &httpHandlerTest{
		http.MethodGet,
		"/database?databaseName=test2",
		nil,
		GetDatabase,
		500,
		nil}
	testHTTPHandler(t, errorNameHandlerTest)
}

func TestCreateOrUpdateDatabase(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)
	_ = state.New("etcd", etcd.Config{
		Endpoints:   []string{clus.Members[0].GRPCAddr()},
		DialTimeout: 2 * time.Second,
	})
	//create success
	normanDatabase := option.Database{Name: "test", NumOfShard: 1, ReplicaFactor: 1}
	normalHandlerTest := &httpHandlerTest{
		http.MethodPost,
		"/database",
		normanDatabase,
		CreateOrUpdateDatabase,
		204,
		nil}
	testHTTPHandler(t, normalHandlerTest)
	// request params error
	noShardDatabase := option.Database{Name: "test", NumOfShard: 0, ReplicaFactor: 1}
	noShardHandlerTest := &httpHandlerTest{
		http.MethodPost,
		"/database",
		noShardDatabase,
		CreateOrUpdateDatabase,
		500,
		nil}
	testHTTPHandler(t, noShardHandlerTest)

	noReplicaDatabase := option.Database{Name: "test", NumOfShard: 1, ReplicaFactor: 0}
	noReplicaHandlerTest := &httpHandlerTest{
		http.MethodPost,
		"/database",
		noReplicaDatabase,
		CreateOrUpdateDatabase,
		500,
		nil}
	testHTTPHandler(t, noReplicaHandlerTest)

	noNameDatabase := option.Database{Name: "", NumOfShard: 1, ReplicaFactor: 1}
	noNameHandlerTest := &httpHandlerTest{
		http.MethodPost,
		"/database",
		noNameDatabase,
		CreateOrUpdateDatabase,
		500,
		nil}
	testHTTPHandler(t, noNameHandlerTest)

}
