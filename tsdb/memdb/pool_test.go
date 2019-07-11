package memdb

import "testing"

func Test_pool(t *testing.T) {
	item1 := tsStoresListPool.get(10)
	tsStoresListPool.put(item1)
	item2 := tsStoresListPool.get(5)
	tsStoresListPool.put(item2)
	tsStoresListPool.get(15)
	tsStoresListPool.get(1)

	item3 := metricStoresListPool.get(10)
	metricStoresListPool.put(item3)
	item4 := metricStoresListPool.get(5)
	metricStoresListPool.put(item4)
	metricStoresListPool.get(15)
	metricStoresListPool.get(1)

}
