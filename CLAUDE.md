# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**LinDB** is a high-performance distributed time series database (TSDB) written in Go. It supports metrics, logs, and traces as data types, uses a custom LinQL (SQL-like) query language parsed via ANTLR4, and integrates Apache Arrow for columnar data processing.

Module path: `github.com/lindb/lindb`

## Commands

### Build
```bash
make build          # Build binaries for darwin/linux/windows (outputs to bin/)
make build-all      # Build with web console included
make run            # Run local standalone mode for dev/debug
```

### Test
```bash
make test-without-lint          # Run all unit tests with race detection + coverage
go test -v ./path/to/package    # Run tests in a specific package
go test -v ./path/to/package -run TestFunctionName  # Run a single test
make e2e-test                   # Run integration tests in e2e/
```

Tests run with `GIN_MODE=release` and `LOG_LEVEL=fatal` to suppress output.

### Lint & Format
```bash
make lint     # Run golangci-lint (installs if missing)
make format   # go fmt ./...
make import   # goimports optimization
make header   # Check/add Apache 2.0 license headers
```

### Code Generation
```bash
make gomock          # Regenerate all *_mock.go files via mockgen
make generate        # Regenerate protobuf/flatbuffers (requires flatbuffers)
make gen-sql-grammar # Regenerate ANTLR4 SQL parser (requires antlr4 v4.13.2)
```

SQL grammar regeneration:
```bash
antlr4 -v 4.13.2 -Dlanguage=Go -listener -visitor -package grammar ./sql/grammar/SQLLexer.g4
antlr4 -v 4.13.2 -Dlanguage=Go -listener -visitor -package grammar ./sql/grammar/SQLParser.g4
```

### Cleanup
```bash
make clean       # Remove mocks, tmp files, binaries
make clean-mock  # Remove all *_mock.go files
make deps        # go mod verify + tidy
```

## Code Generation Workflow

Do **not** run lint or code review after routine edits. `make lint` / `make format` / `make import` and the `code-reviewer` subagent run **only at commit time**, so the same checks are not repeated on every edit.

Rules:
1. After creating or modifying source files, just finish — do **not** run `make lint`, `make format`, `make import`, or the `code-reviewer` subagent.
2. Only when the user asks to commit: run lint (`make lint`, and `make format`/`make import` as needed), fix any errors, run the `code-reviewer` subagent, then commit.
3. Never commit with unresolved lint errors.
4. Running tests to verify correctness (`make test-without-lint`, `go test ...`) is a separate concern and is still fine on demand.

## Architecture

### Node Roles
LinDB runs in three roles. The `app/` directory contains the runtime for each:
- **Broker** (`app/broker/`) — Stateless query/ingestion gateway; routes writes to storage, executes distributed queries
- **Storage** (`app/storage/`) — Stateful storage nodes managing shards and databases
- **Root** (`app/root/`) — Cluster-level management node
- **Standalone** (`app/standalone/`) — Single-process mode embedding all roles (used for development)

Cluster coordination uses **etcd** via `coordinator/` for leader election, service discovery, and state machines.

### SQL Engine (`sql/`)
Query processing pipeline:
1. **Grammar** (`sql/grammar/`) — ANTLR4-generated lexer/parser for LinQL; `.g4` files are the source of truth
2. **Tree** (`sql/tree/`) — AST node definitions and `AstVisitor` that converts ANTLR parse tree to typed AST
3. **Analyzer** (`sql/analyzer/`) — Semantic validation and type resolution
4. **Planner** (`sql/planner/`) — Logical → physical → execution plan generation
5. **Execution** (`sql/execution/`) — DDL, DML, and operator execution

When modifying the SQL grammar, always regenerate the parser files with `make gen-sql-grammar` and commit them alongside the `.g4` changes.

### Storage Engine (`storage/`)
Three data type stores, all under `storage/`:
- **Metric** (`storage/metric/`) — Traditional TSDB metrics; uses custom LSM-tree (`kv/`)
- **Log** (`storage/log/`) — Log data; segments backed by Apache Arrow record batches, uses `github.com/lindb/arrow/pkg/logs`
- **Trace** (`storage/trace/`) — Distributed traces; Arrow-based writer in `app/broker/write/traces/`

Write path: `Engine → Shard → MemoryDatabase → Flusher → KV Store`

### KV Layer (`kv/`)
Custom LSM-tree implementation used by the metric and index stores. Not used by log/trace storage (Arrow-based).

### Index Layer (`index/`)
Inverted index for series metadata. Manages metric names, tag keys/values, field metadata, and series IDs for fast metric lookups.

### Flow Layer (`flow/`)
Orchestrates distributed query execution: node selection, result grouping, filtering, and aggregation fan-out across storage nodes.

### Streaming/CEP Engine (`streaming/`)
Complex Event Processing runtime. Jobs are defined in LinQL, deployed via `streaming/cep/runtime/`, and execute queries against incoming event streams with source/sink connectors.

### Apache Arrow Usage
Arrow (`github.com/apache/arrow-go/v18` and fork `github.com/lindb/arrow`) is used for:
- Log segment storage and scanning (`storage/log/scanner.go`)
- Trace writing (`app/broker/write/traces/`)
- SQL engine execution over record batches (`sql/`)
- The `spi/scalar/` package handles scalar type mapping between LinDB types and Arrow types

### SPI Layer (`spi/`)
Service Provider Interfaces defining the contract between the SQL engine and storage:
- `spi/scalar/` — Scalar type system (maps to Arrow data types)
- `spi/` — Table and column metadata interfaces used by the planner

### Ingestion (`ingestion/`)
Supports multiple write protocols: InfluxDB line protocol, flat file format, and protobuf. Entry point is the broker's HTTP handler.

## Key Conventions

- All source files require the Apache 2.0 license header (enforced by `make header`)
- Mock files are named `*_mock.go` and generated by `mockgen`; never edit them manually
- The `pkg/` directory contains reusable utilities with no business logic dependencies
- `models/` contains shared data structures used across layers; avoid circular imports by keeping models dependency-free
- Build tags: `grocksdb_clean_link` is required when running the standalone server locally (`make run`)
