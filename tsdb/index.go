package tsdb

//go:generate mockgen -source ./index.go -destination=./index_mock.go -package tsdb

type Index interface {
}
