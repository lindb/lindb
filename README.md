<p align="left">
    <img width="400" src="https://github.com/eleme/lindb/wiki/images/readme/lindb_logo.png">
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/lindb/lindb)](https://goreportcard.com/report/github.com/lindb/lindb)
[![Travis CI Build Status](https://travis-ci.org/eleme/lindb.svg?branch=develop)](https://travis-ci.org/eleme/lindb)
[![codecov](https://codecov.io/gh/stone1100/lindb/branch/develop/graph/badge.svg)](https://codecov.io/gh/stone1100/lindb)
[![contribution](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

## What is LinDB?

LinDB("Lin" stands for wisdom in Shanghainese) is a open-source Time Series Database which provides high performance, high availability and horizontal scalability. 

LinDB stores all monitoring data of ELEME Inc, there is 88TB incremental writes per day and 2.7PB total  raw data. 

+ __High performance__

  LinDB takes a lot of best practice of TSDB and implements some optimizations based on the characteristics of time series data. 
  Unlike writing a lot of Continuous-Query for InfluxDB, LinDB supports rollup in specific interval automatically after creating the database. 
  Moreover, LinDB is extremely fast for parallel querying and computing of distributed time series data.

+ __Multi-Active IDCs native__

  LinDB is designed to work under a Multi-Active IDCs cloud architecture. The compute layer of LinDB, called brokers, supports efficient Multi-IDCs aggregation query.

+ __High availability__

  LinDB uses the ETCD cluster to ensure the meta-data is highly available and safely stored. 
  In the event of failure, the Multi-channel replication protocol of WAL will avoid the problem of data inconsistency:  

  1). Only one person in each replication channel is responsible for the authority of the data, so the conflicts will not happen;  

  2). Data reliability is guaranteed: as long as the data that has not been copied in the old leader is not lost, it will be copied to other replication while the old leader is online again; 

+ __Horizontal scalability__

  Tags based sharding strategy in LinDB solves the hotspots problem, and is truly horizontally expanded available by simply adding new broker and storage nodes.
  
+ __Governance capability of metrics__

  To ensure the robustness of the system, LinDB do not assume that users has understood the best practices of using metrics, therefore, LinDB provides the ability of restricting unfriendly user based on metric granularity and tags granularity.

## State of this project

The current develop branch is unstable and is not recommended for production use. LinDB 0.1(what will be the first release version) is currently in the development stage. 
Additional features will arrive during July, we will translate the JAVA version of LinDB currently used under the production environment to Golang as soon as possible.

Once we implement the final feature and replace the LinDB under production environment with the Golang version, LinDB 0.1.0 will be released. At that point, we will move into the stable phase, our intention is to avoid breaking changes to the API and storage file format.

## Build(TODO)

## Architecture

![architecture](https://github.com/lindb/lindb/wiki/images/readme/lindb_architecture.jpg)

## Screenshots
![home](https://github.com/lindb/lindb/wiki/images/readme/home.png)
![search](https://github.com/lindb/lindb/wiki/images/readme/search.png)
![explain](https://github.com/lindb/lindb/wiki/images/readme/explain.png)
![database](https://github.com/lindb/lindb/wiki/images/readme/database.png)

## Contributing

Contributions are welcomed and greatly appreciated. See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

#### CI 
Pull requests should be appropriately labeled, and linked to any relevant [bug or feature tracking issues](https://github.com/lindb/lindb/issues). 
All pull requests will run through Travis-CI. Community contributors should be able to see the outcome of this process by looking at the checks on their PR and fix the build errors.

#### Static Analysis 
This project uses the following linters. Failure during the running of any of these tools results in a failed build. Generally, code must be adjusted to satisfy these tools.

- [gofmt](https://golang.org/cmd/gofmt/) - Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification;
- [golint](https://github.com/golang/lint) - Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes;
- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) - Goimports does everything that gofmt does. Additionally it checks unused imports;
- [errcheck](https://github.com/kisielk/errcheck) - Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases;
- [gocyclo](https://github.com/alecthomas/gocyclo) - Computes and checks the cyclomatic complexity of functions;
- [maligned](https://github.com/mdempsky/maligned) - Tool to detect Go structs that would take less memory if their fields were sorted;
- [dupl](https://github.com/mibk/dupl) - Tool for code clone detection;
- [goconst](https://github.com/jgautheron/goconst) - Finds repeated strings that could be replaced by a constant;
- [gocritic](https://github.com/go-critic/go-critic) - The most opinionated Go source code linter;

## License

LinDB is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.


