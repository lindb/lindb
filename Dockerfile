# Build the manager binary
FROM golang:1.19 as builder
ARG TARGETOS
ARG TARGETARCH
ARG LD_FLAGS


WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY embed.go embed.go
COPY aggregation/ aggregation/
COPY app/ app/
COPY bin/  bin/
COPY cmd/ cmd/
COPY config/ config/
COPY constants/ constants/
COPY coordinator/ coordinator/
COPY data/ data/
COPY docker/ docker/
COPY flow/ flow/
COPY e2e/ e2e/
COPY ingestion/ ingestion/
COPY internal/ internal/
COPY kv/ kv/
COPY metrics/ metrics/
COPY models/ models/
COPY pkg/ pkg/
COPY proto/ proto/
COPY query/ query/
COPY release/ release/
COPY replica/ replica/
COPY rpc/ rpc/
COPY scripts/ scripts/
COPY series/ series/
COPY sql/ sql/
COPY tsdb/ tsdb/
COPY web/ web/
COPY docs/ docs/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN	CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build "${LD_FLAGS}" -o lind ./cmd/lind


FROM centos:latest
WORKDIR /
COPY --from=builder /workspace/lind /usr/bin/lind
RUN mkdir /lindb

USER 0:0