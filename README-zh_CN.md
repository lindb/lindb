<p align="left">
    <img width="400" src="https://github.com/lindb/lindb/wiki/images/readme/lindb_logo.png">
</p>

[![LICENSE](https://img.shields.io/github/license/stone1100/lindb)](https://github.com/lindb/lindb/blob/develop/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/lindb/lindb)](https://goreportcard.com/report/github.com/lindb/lindb)
[![Github Actions Status](https://github.com/lindb/lindb/workflows/LinDB%20CI/badge.svg)](https://github.com/lindb/lindb/actions?query=workflow%3A%22LinDB+CI%22)
[![Github Actions Status](https://github.com/lindb/lindb/workflows/Forntend%20CI/badge.svg)](https://github.com/lindb/lindb/actions?query=workflow%3A%22Forntend+CI%22)
[![codecov](https://codecov.io/gh/lindb/lindb/branch/develop/graph/badge.svg)](https://codecov.io/gh/lindb/lindb)
[![contribution](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

[English](./README.md) | 简体中文

## 简介

LinDB 是一个高性能、高可用并且具备水平拓展性的开源分布式时序数据库。
LinDB 目前在饿了么存储了全量的监控数据，每天增量写入 88TB，全量数据总计 2.7PB。

+ __高性能__

  LinDB 参考了其他时序数据库里的最佳实践，并且在此基础上做了些优化。
  和 InfluxDB 中的 Continuous-Query 不一样，在创建数据库之后，LinDB支持根据指定的时间间隔自动对数据进行rollup。
  不仅如此，LinDB 在并行查询和计算分布式时序数据方面速度相当快。

+ __原生的多活支持__

  LinDB 在设计阶段就是为了在饿了么的多活系统架构下工作的。LinDB的计算层（brokers）支持高效的跨机房聚合查询。

+ __高可用__

  LinDB 使用 ETCD 集群以保证原数据的高可用和存储安全（也支持单机模式启动）。
  当出现故障时，多通道复制的WAL协议可以避免数据的不一致情况：

  1). 任一复制通道只有一名成员负责数据的权威性，因此避免了数据冲突;

  2). 数据可靠性保障：只要先前的leader数据并未丢失，当其重新上线时，这些数据便将被复制到其他的副本中；

+ __水平拓展性__

  LinDB 使用了基于 Series(Tags) 的 sharding 策略，可以真正解决热点问题（另外也考验了系统设计的能力），简单的增加新的 broker 和 storage 节点即可实现水平拓展。

+ __Metric的治理能力__

  为了保证系统的健壮性，LinDB 并不假设所有用户都已了解 metric 使用的最佳实践。因此，LinDB 提供了 metric 和 tag 粒度的限制能力，避免整体集群在部分用户错误的使用姿势下崩溃。

## 项目状态

当前的分支还不稳定，并且不推荐生产环境的使用，LinDB 0.1 还处在开发阶段。新的功能点会在七月、八月陆续完成开发。
目前我们在饿了么生产环境中使用的是 JAVA 版本的 LinDB，开源的 Go 版本不仅仅是对 JAVA 版的简单翻译，在很多地方都进行了重新设计，同时受开发资源的影响，功能开发速度较慢。
当完成基本的功能开发后，我们会 release 0.1.0 的 alpha 版本。在此之后，当进入稳定版的迭代过程时，我们会尽量避免 API 的不兼容更新以及存储文件格式的变更。

#### JAVA 版 LinDB 相关文章
- [数据分析之时序数据库](https://zhuanlan.zhihu.com/p/36804890)
- [分布式时序数据库 - LinDB](https://zhuanlan.zhihu.com/p/35998778)

## 编译

### 依赖

在本地编译 LinDB 需要以下工具：
- [Go >=1.14](https://golang.org/doc/install)
- [Make tool](https://www.gnu.org/software/make/)
- [Yarn](https://classic.yarnpkg.com/en/docs/install)

### 环境设置

设置 Go path 与 PATH 环境变量，GOPATH 通常默认为 `$HOME/go`

```
# 将以下命令添加到 ~/.bashrc 或 ~/.bash_profile 文件中
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# 使用 source 命令使 bashrc 文件立刻生效
$ source ~/.bashrc
```

### 获取代码

```
git clone https://github.com/lindb/lindb.git
cd lindb
```

### 编译源代码

仅编译 LinDB 后端（无管理界面）

```
make build
```

同时编译 LinDB 前端与后端

```
make build-all
```

### 测试

```
make test
```

### 管理界面(开发者)

启动 LinDB 前端应用
```
cd web
yarn start
```

可以通过  [localhost port 3000](http://localhost:3000/) 来访问

## 部署

LinDB 有两种部署方式：单机模式与集群模式

### 单机模式

单机模式包含了内嵌的 broker, storage 和 etcd 模块

```
./bin/lind standalone init-config
./bin/lind standalone run
```
在使用 `make build-all` 编译后，前端文件会被打包在 broker 模块中
可以通过  [localhost port 9000](http://localhost:9000/) 来访问控制台

### 集群模式 (todo)


## 架构

![architecture](https://github.com/lindb/lindb/wiki/images/readme/lindb_architecture.jpg)

## 贡献代码

我们非常期待有社区爱好者能加入我们一起参与开发，[CONTRIBUTING](CONTRIBUTING.md) 是一些简单的 PR 的规范，对于 一个 PR 中的多个 commit，我们会根据情况在合并时做 squash 并进行归类，以方便后续查看回溯。

#### CI
PR 应当带上合适的标签，并且关联到已有的 issue 上 [issues](https://github.com/lindb/lindb/issues)。
所有的 PR 都会在 GITHUB-Actions 进行测试，社区贡献者需要关注 CI 的结果，对未通过的错误进行修复。

#### 静态检查
我们使用了以下的检查器，所有代码都需要针对以下工具做一些调整。

- [gofmt](https://golang.org/cmd/gofmt/) - Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification;
- [golint](https://github.com/golang/lint) - Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes;
- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) - Goimports does everything that gofmt does. Additionally it checks unused imports;
- [errcheck](https://github.com/kisielk/errcheck) - Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases;
- [gocyclo](https://github.com/alecthomas/gocyclo) - Computes and checks the cyclomatic complexity of functions;
- [maligned](https://github.com/mdempsky/maligned) - Tool to detect Go structs that would take less memory if their fields were sorted;
- [dupl](https://github.com/mibk/dupl) - Tool for code clone detection;
- [goconst](https://github.com/jgautheron/goconst) - Finds repeated strings that could be replaced by a constant;
- [gocritic](https://github.com/go-critic/go-critic) - The most opinionated Go source code linter;

## 开源许可协议

LinDB 使用 Apache 2.0 协议， 点击 [LICENSE](LICENSE) 查看详情。


