package operator

// A distributed SQL query engine designed for time series(ref:https://trino.io/).
// It supports a variety of SQL operations, similar to traditional SQL databases,
// but also includes optimizations and features suited for time series environments.
// Here are some key operators and concepts in LinDB:
//
// 1. TableScan: reads data from a data source(metric/log etc.)
// 2. Projection: applies transformations to the input data, suchs computing expressions based on input columns
// 3. Filter: filters rows that do not meet the condition specified
// 4. Aggregation: performs calculations across a set of rows that are grouped together based on one or more columns(sum/avg/count etc.)
// 5. Join: combines rows from two or more tables based on a related column between them
// 6: Sort: order the data based on one or more columns(asc/desc)
// 7: Limit: limits the number of rows returned by a query
// 8: Window Function?
// 9: TopN?
// 10: Hash?
