package option

import "github.com/eleme/lindb/pkg/interval"

type ShardOption struct {
	Behind       int64         // allowed timestamp write behind
	Ahead        int64         // allowed timestamp write ahead
	Interval     int64         // interval
	IntervalType interval.Type // interval type
}
