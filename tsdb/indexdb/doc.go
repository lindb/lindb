package indexdb

/*

Metric id mapping storage format using bbolt:
+-------------+
| root Bucket |
+------+------+
       |       +---------------+
       +------>+ series bucket +
               +-------+-------+
                       |         +-----------------+
                       +-------->+ metric bucket 1 +-----> series id sequence of metric level
                       |         +--------+--------+
                       |                  |         +------------------------+
                       |                  +-------->+ tags hash 1->series id |
                       |                  |         +------------------------+
                       |                  |         +------------------------+
                       |                  +-------->+ tags hash n->series id |
                       |                            +------------------------+
                       |         +-----------------+
                       +-------->+ metric bucket n |
                                 +-----------------+
*/
