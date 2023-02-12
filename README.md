<p align="left">
    <img width="400" src="https://github.com/lindb/lindb/wiki/images/readme/lindb_logo.png">
</p>

[![LICENSE](https://img.shields.io/github/license/lindb/lindb)](https://github.com/lindb/lindb/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/lindb/lindb)](https://goreportcard.com/report/github.com/lindb/lindb)
[![Github Actions Status](https://github.com/lindb/lindb/workflows/LinDB%20CI/badge.svg)](https://github.com/lindb/lindb/actions?query=workflow%3A%22LinDB+CI%22)
[![Github Actions Status](https://github.com/lindb/lindb/workflows/Forntend%20CI/badge.svg)](https://github.com/lindb/lindb/actions?query=workflow%3A%22Forntend+CI%22)
[![codecov](https://codecov.io/gh/lindb/lindb/branch/main/graph/badge.svg)](https://codecov.io/gh/lindb/lindb)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/lindb/lindb)
[![contribution](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

English | [简体中文](./README-zh_CN.md)

## What is LinDB?

LinDB is an open-source Time Series Database which provides high performance, high availability and horizontal scalability. 

- [Key features](https://lindb.io/guide/introduction.html#key-features)
- [Guide](https://lindb.io/guide/introduction.html)
- [Quick start](https://lindb.io/guide/get-started.html)
- [Design](https://lindb.io/design/architecture.html)
- [Architecture](#architecture)
- [Admin UI](#admin-ui)

## Build

### Prerequisites

To build LinDB from source you require the following on your system.

- [Go >=1.16](https://golang.org/doc/install)
- [Make tool](https://www.gnu.org/software/make/)
- [Yarn](https://classic.yarnpkg.com/en/docs/install)

### Get the code

```
git clone https://github.com/lindb/lindb.git
cd lindb
```

### Build from source

To build only LinDB core.(without web console)

```
make build
```

To build both LinDB core and frontend.

```
make build-all
```

### Test

```
make test
```

### Access web interface(for developer)

Start the node.js app to view LinDB web interface in dev mode.

```
cd web
yarn install 
yarn dev
```

You can access the LinDB web interface on your [localhost port 3000](http://localhost:3000/)

## Architecture

![architecture](./docs/images/architecture.png)

## Admin UI

Some admin ui snapshots.

### Overview

![overview](./docs/images/overview.png)

### Monitoring Dashboard

![dashboard](./docs/images/dashboard.png)

### Replication State

![replication](./docs/images/replication_shards.png)

### Data Explore

![explore](./docs/images/data_explore.png)

### Explain

![explain](./docs/images/data_search_explain.png)

## Contributing

Contributions are welcomed and greatly appreciated. See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

#### CI 
Pull requests should be appropriately labeled, and linked to any relevant [bug or feature tracking issues](https://github.com/lindb/lindb/issues). 
All pull requests will run through GITHUB-Actions. Community contributors should be able to see the outcome of this process by looking at the checks on their PR and fix the build errors.

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
