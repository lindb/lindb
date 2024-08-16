## [v0.4.1](https://github.com/lindb/lindb/releases/tag/v0.4.1) - 2024-08-16

See [code changes](https://github.com/lindb/lindb/compare/v0.4.0...v0.4.1).

### 🚀 New features

- [enhance]: remove installation guides from release notes by @stone1100 in #1043

### 🐛 Bug fixes

- [bug]: fix cannot flush metric data when server shutdown by @stone1100 in #1044


## [v0.4.0](https://github.com/lindb/lindb/releases/tag/v0.4.0) - 2024-07-22

See [code changes](https://github.com/lindb/lindb/compare/v0.3.1...v0.4.0).

### 🚀 New features

- [opt]: reduce memory database memory usage by @stone1100 in #1030
- [opt]: ignore histogram bucket if count<0 by @stone1100 in #1033
- [opt]: opt metric field memory store by @stone1100 in #1035
- [feat]: memory database approximate memory size by @stone1100 in #1036
- [feat]: cleanup memory metric meta/index if not used by @stone1100 in #1037
- [enhance]: add self-monitoring metric/docs by @stone1100 in #1040
- [feat]: add tsdb limits(namespace/metric/field/tag key/series) by @stone1100 in #1041

### 🐛 Bug fixes

- [bug]: fix build docker fail by @stone1100 in #1031
- [bug]: miss makezero in slice init by @stone1100 in #1032
- [opt]: ignore histogram bucket if count<0 by @stone1100 in #1033
- [bug]: fix get wrong data from memory database by @stone1100 in #1034
- [bug]: fix wal ack invalid seq msg by @stone1100 in #1038

### 💬 Others

- [docs]: update architecture image by @stone1100 in #985
- [refactor]: use common lib(logger/timeutil/fileutil) by @stone1100 in #986
- [feat]: remove duplicate agg types by @joyant in #989
- [opt]: succinct trie by @stone1100 in #990
- [bug]: fix nil metadata and refactor lindcli by @joyant in #993
- [bug:#994]: fix not in statement returns an empty result by @joyant in #996
- [feat]: build index on LinDB common kv store by @stone1100 in #997
- [chore]: upgrade go version for ci by @stone1100 in #998
- [feat:#912]: support comparison binary operations by @joyant in #1000
- [feat:#995]: create database by with statement by @joyant in #1002
- [bug]: fix read index data panic after kv store compact by @stone1100 in #1004
- [chore]: disable golangci-lint cache by @stone1100 in #1006
- [feat]: memory database estimate heap size by @stone1100 in #1005
- [refactor]: memory database data loader by @stone1100 in #1009
- [refactor]: only one storage cluster under broker cluster by @kevin6025 in #1008
- [bug:#1012]: fix show metrics returns incorrect results by @joyant in #1013
- [feat:#1001]: support promql by @joyant in #1014
- [enhance:#1007]: reduce goroutine when write too many data families by @joyant in #1016
- [chore]: rebase v0.3.0_bug_fix by @stone1100 in #1025
- [bug:#1019]: fix storage node goes dead when gc pause by @stone1100 in #1026
- [opt]: compress write data by re-use snappy streaming by @stone1100 in #1027
- [docs]: add Japanese README file by @eltociear in #1039


## [v0.3.1](https://github.com/lindb/lindb/releases/tag/v0.3.1) - 2024-04-21

See [code changes](https://github.com/lindb/lindb/compare/v0.3.0...v0.3.1).

### 🐛 Bug fixes

- [bug]: fix storage panic when write old data point by @stone1100 in #1020

## [v0.3.0](https://github.com/lindb/lindb/releases/tag/v0.3.0) - 2023-08-29

See [code changes](https://github.com/lindb/lindb/compare/v0.2.6...v0.3.0).

### 🚀 New features

- [chore]: add style lint by @stone1100 in #977
- [enhance]: add iconfont by @stone1100 in #979
- [feat]: CI support arm64 images by @dongjiang1989 in #981

### 🐛 Bug fixes

- [bug]: fix theme color palette setting by @stone1100 in #978
- [bug]: use metric level namespace if set by @stone1100 in #980
- [bug]: fix data explore not support namespace by @stone1100 in #982
- [bug]: fix diff namespace conflict by @stone1100 in #983


## [v0.2.6](https://github.com/lindb/lindb/releases/tag/v0.2.6) - 2023-04-23

See [code changes](https://github.com/lindb/lindb/compare/v0.2.5...v0.2.6).

### 🚀 New features

- [feat]: support select *(query all fields) by @stone1100 in #974
- [chore]: build darwin arm64 package by @stone1100 in #975


## [v0.2.5](https://github.com/lindb/lindb/releases/tag/v0.2.5) - 2023-04-17

See [code changes](https://github.com/lindb/lindb/compare/v0.2.4...v0.2.5).

### 🐛 Bug fixes

- [bug]: darwin package no cpu stats data by @stone1100 in #968
- [bug]: lost web console when build package via github action by @stone1100 in #969
- [bug]: fix build darwin package fail by @stone1100 in #971


## [v0.2.4](https://github.com/lindb/lindb/releases/tag/v0.2.4) - 2023-04-16

See [code changes](https://github.com/lindb/lindb/compare/v0.2.3...v0.2.4).

### 🚀 New features

- [enhance]: consume group wait strategy when WAL is empty by @stone1100 in #960

### 💬 Others

- [chore]: use golangci/golangci-lint-action by @damnever in #956
- [chore]: add twitter badge to readme by @stone1100 in #958


## [v0.2.3](https://github.com/lindb/lindb/releases/tag/v0.2.3) - 2023-04-03

See [code changes](https://github.com/lindb/lindb/compare/v0.2.2...v0.2.3).

### Added features

- docker hub sync github actions

### Optimized

- join file paths

### Fixed

- fix get diff agg result using diff group by interval

## [v0.2.2](https://github.com/lindb/lindb/releases/tag/v0.2.2) - 2023-03-28

See [code changes](https://github.com/lindb/lindb/compare/v0.2.1...v0.2.2).

### Added features

- support auto fill group by time interval based on query range  
- read config from env
- print config value when server startup
- support release with a docker image

### Optimized

- print error message when handle http request failure
- add cli tools to docker image
- need to check rollup Interval when memory database refreshing
- check database's intervals if valid

### Fixed

- fix time picker no display last 30 days
- cannot create database if storage not exist
- got unexpected result when async load family data
- fix cannot rollup data for windows
- fix miss some source family data when rollup job

## New Contributors

- @Sn0rt made their first contribution in https://github.com/lindb/lindb/pull/927

## [v0.2.1](https://github.com/lindb/lindb/releases/tag/v0.2.1) - 2023-03-04

See [code changes](https://github.com/lindb/lindb/compare/v0.2.0...v0.2.1).

### Added features

- database level read/write limit configure
- slow sql

### Optimized

- metric chart show error message if query fail

### Fixed

- fix nil point when parse wrong sql

## [v0.2.0]https://github.com/lindb/lindb/releases/tag/v0.2.0) - 2023-01-29

See [code changes](https://github.com/lindb/lindb/compare/v0.1.1...v0.2.0).

### Added features

- cross multiple IDCs query;
- configure logic database/broker cluster;
- add root server for cross IDCs;
- explore root's coordinate metadata;
- root self-monitoring;
- view root state;

### Optimized

- view admin UI based env vars;
- de-register broker node when shutdown;
- stop write ahead log too slow;
- refactor query stats metric;
- refactor explain query stats;

### Fixed

- fix directory traversal security issue;
- fix alive/create/expire task metric no data;

## [v0.1.1](https://github.com/lindb/lindb/releases/tag/v0.1.1) - 2022-12-04

See [code changes](https://github.com/lindb/lindb/compare/v0.1.0...v0.1.1).

### Added features

- web admin console support zh_CN;
- add detail doc link for each page;
- recover database metadata from local storage if etcd data lost;

### Optimized

- print version format;
- lind-cli logmsg via file;
- chart's y-axes start with 0 if min value <=0;

### Fixed

- fix delete active wal when not set write ahead time;

## [v0.1.0]https://github.com/lindb/lindb/releases/tag/v0.1.0) - 2022-11-22

### Added features

- Metadata coordinator(database/broker/storage/master etc.);
- Distributed query engine, SQL supported;
- Data supports distributed storage;
- Supports write ahead log;
