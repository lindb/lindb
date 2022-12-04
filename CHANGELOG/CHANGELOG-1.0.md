# Change logs

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

