<p align="left">
    <img width="400" src="https://github.com/lindb/lindb/wiki/images/readme/lindb_logo.png">
</p>

[![LICENSE](https://img.shields.io/github/license/lindb/lindb)](https://github.com/lindb/lindb/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/lindb/lindb)](https://goreportcard.com/report/github.com/lindb/lindb)
[![LinDB CI](https://github.com/lindb/lindb/actions/workflows/lind.yml/badge.svg)](https://github.com/lindb/lindb/actions/workflows/lind.yml)
[![Frontend CI](https://github.com/lindb/lindb/actions/workflows/frontend.yml/badge.svg)](https://github.com/lindb/lindb/actions/workflows/frontend.yml)
[![Docker Latest](https://github.com/lindb/lindb/actions/workflows/docker-latest.yml/badge.svg)](https://github.com/lindb/lindb/actions/workflows/docker-latest.yml)
[![Docker Release](https://github.com/lindb/lindb/actions/workflows/docker-release.yml/badge.svg)](https://github.com/lindb/lindb/actions/workflows/docker-release.yml)
[![codecov](https://codecov.io/gh/lindb/lindb/branch/main/graph/badge.svg)](https://codecov.io/gh/lindb/lindb)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/lindb/lindb)
[![contribution](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)
[![Twitter Follow](https://img.shields.io/twitter/follow/lindb_io?style=social)](https://twitter.com/intent/follow?screen_name=lindb_io)

[English](./README.md) | [简体中文](./README-zh_CN.md) | 日本語

## LinDBとは？

LinDBは、高性能、高可用性、水平スケーラビリティを提供するオープンソースの時系列データベースです。

- [主な特徴](https://lindb.io/guide/introduction.html#key-features)
- [ガイド](https://lindb.io/guide/introduction.html)
- [クイックスタート](https://lindb.io/guide/get-started.html)
- [設計](https://lindb.io/design/architecture.html)
- [アーキテクチャ](#アーキテクチャ)
- [管理UI](#管理-ui)

## ビルド

### 前提条件

LinDBをソースからビルドするには、以下のツールがシステムに必要です。

- [Go >=1.21](https://golang.org/doc/install)
- [Make tool](https://www.gnu.org/software/make/)
- [Yarn](https://classic.yarnpkg.com/en/docs/install)

### コードの取得

```
git clone https://github.com/lindb/lindb.git
cd lindb
```

### ソースからのビルド

LinDBコアのみをビルドします。（Webコンソールなし）

```
make build
```

LinDBコアとフロントエンドの両方をビルドします。

```
make build-all
```

### テスト

```
make test
```

### Webインターフェースへのアクセス（開発者向け）

LinDBのWebインターフェースを開発モードで表示するために、Node.jsアプリを起動します。

```
cd web
yarn install 
yarn dev
```

LinDBのWebインターフェースには、[localhostのポート3000](http://localhost:3000/)でアクセスできます。

## アーキテクチャ

![architecture](./docs/images/architecture.png)

## 管理UI

管理UIのスナップショットの一部です。

### 概要

![overview](./docs/images/overview.png)

### 監視ダッシュボード

![dashboard](./docs/images/dashboard.png)

### レプリケーション状態

![replication](./docs/images/replication_shards.png)

### データ探索

![explore](./docs/images/data_explore.png)

### 説明

![explain](./docs/images/data_search_explain.png)

## コントリビューション

コントリビューションは歓迎され、非常に感謝されます。パッチの提出方法やコントリビューションのワークフローについては、[CONTRIBUTING](CONTRIBUTING.md)をご覧ください。

#### CI 
プルリクエストには適切なラベルを付け、関連する[バグまたは機能追跡の問題](https://github.com/lindb/lindb/issues)にリンクする必要があります。
すべてのプルリクエストはGITHUB-Actionsを通じて実行されます。コミュニティのコントリビューターは、プルリクエストのチェックを見てビルドエラーを修正することで、このプロセスの結果を確認できます。

#### 静的解析 
このプロジェクトでは、以下のリンターを使用しています。これらのツールの実行中に失敗すると、ビルドが失敗します。一般的に、これらのツールを満たすようにコードを調整する必要があります。

- [gofmt](https://golang.org/cmd/gofmt/) - Gofmtはコードがgofmtされているかどうかをチェックします。デフォルトでは、コードの簡略化をチェックするために-sオプションで実行されます。
- [golint](https://github.com/golang/lint) - Golintはgofmtとは異なります。GofmtはGoのソースコードを再フォーマットしますが、golintはスタイルの間違いを出力します。
- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) - Goimportsはgofmtが行うすべてのことを行います。さらに、未使用のインポートをチェックします。
- [errcheck](https://github.com/kisielk/errcheck) - Errcheckは、Goプログラムでチェックされていないエラーをチェックするプログラムです。これらのチェックされていないエラーは、場合によっては重大なバグになる可能性があります。
- [gocyclo](https://github.com/alecthomas/gocyclo) - 関数の循環的複雑度を計算してチェックします。
- [maligned](https://github.com/mdempsky/maligned) - フィールドがソートされている場合にメモリを節約できるGo構造体を検出するツールです。
- [dupl](https://github.com/mibk/dupl) - コードクローン検出ツールです。
- [goconst](https://github.com/jgautheron/goconst) - 定数に置き換えることができる繰り返しの文字列を見つけます。
- [gocritic](https://github.com/go-critic/go-critic) - 最も意見のあるGoソースコードリンターです。

## ライセンス

LinDBはApache 2.0ライセンスの下で提供されています。詳細については、[LICENSE](LICENSE)ファイルを参照してください。
