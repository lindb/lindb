package queue

/**
FanOutQueue structure
+---------------------------------------------------------+
| Tail             FanOutQueue                    Head    |
|                                                         |
|                                                         |
|    +-----------+   +-----------+                        |
|    |  Segment  |   |  Segment..|                        |
|    |  Begin    |   |           |                        |
|    |  End      |   |           |                        |
|    |  Append   |   |           |                        |
|    |  Read     |   |           |                        |
|    |           |   |           |                        |
|    +----+------+   +------+----+                        |
|         ^                 ^                             |
+---------+-----------------+-----------------------------+
|         |                 |                             ^
|         |                 |                             |
+---------+--+   +----------+-+                           | Append
|  FanOut    |   |   FanOut.. |                           |
|  Name      |   |            |                           + +--------+
|  Consume   |   |            |                             | Data...|
|  Get       |   |            |                             |        |
|  Ack       |   |            |                             |        |
|            |   |            |                             +--------+
+------------+   +------------+

Segment structure
+-----------------------------------------------+
|            Segment                            |
|                                               |
|    +----------------------+    +----------+   |
|    |    IndexPage         |    | DataPage |   |
|    |  dataOffset dataLen+------->Message1 |   |
|    |  (4 bytes) (4 bytes) |    | Message2 |   |
|    |  ....                |    | ....     |   |
|    |                      |    |          |   |
|    |                      |    |          |   |
|    +----------------------+    +----------+   |
|                                               |
+-----------------------------------------------+


directory structure
/fanoutQueueDir
	queue.meta
	/segments
		0.idx
		0.dat
		100.idx
		100.dat
		...
	/fanOut
		fanOutName1.meta
		fanOutName2.meta
		...
*/
