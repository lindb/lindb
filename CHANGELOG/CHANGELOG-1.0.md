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
