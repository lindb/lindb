package replication

/**
Channel represents a persistent, replicable storage for a data shard in Broker.
One data shard may be copied to many storage for fault-tolerance.
Replicator handles the details of copping data shard to a target Storage node.

                                                             Storage
                              +--------------+               +-------------------+
                              |Replicator    |               |  StorageService   |
                              |              |  Write        |                   |
Broker                        |Target        +-------------->+  SequenceManager  |
+----------------------+      |Cluster       |  Stream RPC   |  HeadSeq          |
| Channel              |      |Database      <---------------+  AckSeq           |
|                      +------+ShardID       |  Ack Reset    |                   |
| Cluster              |      |Pending       |               |                   |
| Database             |      |Stop          |               |                   |
| ShardID              |      |              |               |                   |
| Write                |      +--------------+               +-------------------+
| GetOrCreateReplicator|
| Targets              |
|                      |
|                      |
+----------------------+      Replicators......

*/
