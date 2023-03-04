# Change logs

## [0.2.1] - 2023-03-04

See [code changes](https://github.com/lindb/lindb/compare/v0.2.0...v0.2.1).

### Added features

- database level read/write limit configure
- slow sql

### Optimized

- metric chart show error message if query fail

### Fixed
- fix nil point when parse wrong sql

## [0.2.0] - 2023-01-29

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

## [0.1.1] - 2022-12-04

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

## [0.1.0] - 2022-11-22

### Added features

- Metadata coordinator(database/broker/storage/master etc.);
- Distributed query engine, SQL supported;
- Data supports distributed storage;
- Supports write ahead log;

