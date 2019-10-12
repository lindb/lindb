package parallel

/*

ref: drill query execution(https://drill.apache.org/docs/drill-query-execution/)


distribution query and compute, like fork-join pattern, and executes all tasks asynchronously
1) sampling aggregation and down sampling without grouping
                        +---------+
                        +  client +
       	                +----+----+
                             |
                             v
                        +----+----+
       +--------------->+   root  +<---------------+
       |                +----+----+          	   |
       |                     |                     |
       |                     |                     |
       v                     v                     v
+------+------+       +------+------+       +------+------+
|storage node1|       |storage node2|       |storage nodeN|
+------+------+       +------+------+       +------+------+

2) complex query with grouping

                       +---------+
                       +  client +
                       +----+----+
                            |
                            v
                       +----+----+
                  +--->+   root  +<---+
                  |    +---------+    |
                  |                   |
                  v					  v
    	   +------+------+     +------+------+
      	   |intermediate1|     |intermediateN|
           +------+------+     +------+------+
                  ^                   ^
       			  |            	      |
                  |                   |
                  v                   v
       +--------->+<-------<+>------->+<----------+
       |               	    |   	              |
       |              	    |   	              |
       v                    v                     v
+------+------+      +------+------+       +------+------+
|storage node1|      |storage node2|       |storage nodeN|
+-------------+      +-------------+       +-------------+
*/
