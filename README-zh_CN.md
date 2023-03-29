<p align="left">
    <img width="400" src="https://github.com/lindb/lindb/wiki/images/readme/lindb_logo.png">
</p>

[![LICENSE](https://img.shields.io/github/license/lindb/lindb)](https://github.com/lindb/lindb/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/lindb/lindb)](https://goreportcard.com/report/github.com/lindb/lindb)
[![Github Actions Status](https://github.com/lindb/lindb/workflows/LinDB%20CI/badge.svg)](https://github.com/lindb/lindb/actions?query=workflow%3A%22LinDB+CI%22)
[![Github Actions Status](https://github.com/lindb/lindb/workflows/Forntend%20CI/badge.svg)](https://github.com/lindb/lindb/actions?query=workflow%3A%22Forntend+CI%22)
[![Github Actions Status](https://github.com/lindb/lindb/workflows/Docker%20Latest/badge.svg)](https://github.com/lindb/lindb/actions?query=workflow%3A%22Docker+Latest%22)
[![Github Actions Status](https://github.com/lindb/lindb/workflows/Docker%20Release/badge.svg)](https://github.com/lindb/lindb/actions?query=workflow%3A%22Docker+Release%22)
[![codecov](https://codecov.io/gh/lindb/lindb/branch/main/graph/badge.svg)](https://codecov.io/gh/lindb/lindb)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/lindb/lindb)
[![contribution](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

[English](./README.md) | 简体中文

## 简介

LinDB 是一个高性能、高可用并且具备水平拓展性的开源分布式时序数据库。

- [主要特性](https://lindb.io/zh/guide/introduction.html#主要特性)
- [用户指南](https://lindb.io/zh/guide/introduction.html)
- [快速开始](https://lindb.io/zh/guide/get-started.html)
- [设计](https://lindb.io/zh/design/architecture.html)
- [架构](#架构)
- [Admin UI](#admin-ui)

## 编译

### 依赖

在本地编译 LinDB 需要以下工具：
- [Go >=1.16](https://golang.org/doc/install)
- [Make tool](https://www.gnu.org/software/make/)
- [Yarn](https://classic.yarnpkg.com/en/docs/install)

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
yarn install
yarn dev
```

可以通过  [localhost port 3000](http://localhost:3000/) 来访问

## 架构

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

## JAVA 版 LinDB 相关文章
- [数据分析之时序数据库](https://zhuanlan.zhihu.com/p/36804890)
- [分布式时序数据库 - LinDB](https://zhuanlan.zhihu.com/p/35998778)

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


