package concurrent

// PoolStat is the statistics data for pool
type PoolStat struct {
	AliveWorkers   int
	CreatedWorkers int
	KilledWorkers  int
	ConsumedTasks  int
}
