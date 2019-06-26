package option

// Database define database config
type Database struct {
	Name          string `json:"name"`
	NumOfShard    int    `json:"numOfShard"`
	ReplicaFactor int    `json:"replicaFactor"`
}
