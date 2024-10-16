package infoschema

import (
	"context"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	protoMetaV1 "github.com/lindb/lindb/proto/gen/v1/meta"
)

func (r *reader) suggestNamespaces(database, ns string, limit int64) ([]string, error) {
	partitions, err := r.metadataMgr.GetPartitions(database, "", "")
	if err != nil {
		return nil, err
	}
	var values []string
	for node := range partitions {
		conn, err := grpc.Dial(node.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := protoMetaV1.NewMetaServiceClient(conn)
		resp, err := client.SuggestNamespace(context.TODO(), &protoMetaV1.SuggestRequest{
			Database:  database,
			Namespace: ns,
			Limit:     limit,
		})
		if err != nil {
			return nil, err
		}
		values = append(values, resp.Values...)
	}
	return values, nil
}

func (r *reader) suggestTables(database, ns, table string, limit int64) ([]string, error) {
	partitions, err := r.metadataMgr.GetPartitions(database, "", "")
	if err != nil {
		return nil, err
	}
	var values []string
	for node := range partitions {
		conn, err := grpc.Dial(node.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := protoMetaV1.NewMetaServiceClient(conn)
		resp, err := client.SuggestTable(context.TODO(), &protoMetaV1.SuggestRequest{
			Database:  database,
			Namespace: ns,
			Table:     table,
			Limit:     limit,
		})
		if err != nil {
			return nil, err
		}
		values = append(values, resp.Values...)
	}
	return values, nil
}

// getStateFromStorage returns the state from storage cluster.
func (r *reader) getStateFromStorage(path string, params map[string]string, newStateFn func() any) (any, error) {
	liveNodes := r.metadataMgr.GetStorageNodes()
	return r.fetchStateData(
		lo.Map(liveNodes, func(node models.StatefulNode, index int) models.Node { return &node }),
		path, params, newStateFn)
}

// fetchStateData fetches the state metric from each live node.
func (r *reader) fetchStateData(nodes []models.Node, path string, params map[string]string, newStateFn func() any) (any, error) {
	size := len(nodes)
	if size == 0 {
		return nil, nil
	}
	result := make([]any, size)
	var wait sync.WaitGroup
	wait.Add(size)
	for idx := range nodes {
		i := idx
		go func() {
			defer wait.Done()
			node := nodes[i]
			address := node.HTTPAddress()
			state := newStateFn()
			_, err := resty.New().R().SetQueryParams(params).
				SetHeader("Accept", "application/json").
				SetResult(&state).
				Get(address + constants.APIVersion1CliPath + path)
			if err != nil {
				r.logger.Error("get state from storage node", logger.String("url", address), logger.Error(err))
				return
			}
			result[i] = state
		}()
	}
	wait.Wait()
	rs := make(map[string]any)
	for idx := range nodes {
		rs[nodes[idx].Indicator()] = result[idx]
	}
	return rs, nil
}
