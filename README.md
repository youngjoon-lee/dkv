# Distributed Key-value Store

This project is a distributed key-value store that users can store key-value pairs to.
All key-value pairs stored are replicated to follower nodes in case the leader node becomes unavailable.

* [Overview](#overview)
* [Installation](#installation)
* [Getting started](#getting-started)
    * [Running nodes](#running-nodes)
    * [Storing a key-value pair](#storing-a-key-value-pair)
    * [Querying a key-value pair](#querying-a-key-value-pair)
* [Development](#development)
    * [Building the project](#building-the-project)
    * [Development Roadmap](#development-roadmap)
* [License](#license)


## Overview

This distributed key-value store satisfies C (consistency) and P (partition tolerance) among the [CAP theorem](https://en.wikipedia.org/wiki/CAP_theorem).
For consistency, it guarantees [strong consistency](https://en.wikipedia.org/wiki/Strong_consistency).

To replicate operations, this project implements a replication layer from scratch,
but the main idea of the replication is adopted from the [Raft](https://raft.github.io/) consensus algorithm.


## Installation

```bash
make install  # will build and put a binary under your $GOPATH/bin/
```


## Getting started

### Running nodes

A cluster with 3 nodes can be run by the following commands:
```bash
DKV_LOG_LEVEL=info \
DKV_RPC_PORT=7070 \
DKV_REST_PORT=7071 \
DKV_DB_PATH="a.db" \
DKV_NODE_ID=a \
DKV_CLUSTER="a@127.0.0.1:7070,b@127.0.0.1:8080,c@127.0.0.1:9090" \
dkv

DKV_LOG_LEVEL=info \
DKV_RPC_PORT=8080 \
DKV_REST_PORT=8081 \
DKV_DB_PATH="b.db" \
DKV_NODE_ID=b \
DKV_CLUSTER="a@127.0.0.1:7070,b@127.0.0.1:8080,c@127.0.0.1:9090" \
dkv

DKV_LOG_LEVEL=info \
DKV_RPC_PORT=9090 \
DKV_REST_PORT=9091 \
DKV_DB_PATH="c.db" \
DKV_NODE_ID=c \
DKV_CLUSTER="a@127.0.0.1:7070,b@127.0.0.1:8080,c@127.0.0.1:9090" \
dkv
```

A `dkv` process exposes two endpoints, gRPC and REST. The REST endpoint is auto-generated from gRPC by [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

Each node should be configured with a unique ID, `DKV_NODE_ID`, and a corresponding DB path that is used to create/open a [BoltDB](https://github.com/etcd-io/bbolt) instance.

Also, the cluster information should be specified as a `DKV_CLUSTER` environment variable, with a comma-separated string.
The first node specified in the `DKV_CLUSTER` is considered as a leader node.

Currently, dynamic cluster configuration is not supported yet.

### Storing a key-value pair

After all nodes are started, users can send a HTTP POST request to a leader of the cluster to store a key-value pair to the store.
```bash
curl -X POST localhost:7071/v0/kv -d '{"key":"aGVsbG8=","value":"d29ybGQ="}'
```
Please note that the data type of key and value is a base64-encoded string.
These key and value are decoded internally and stored in the DB as a byte array each.

### Querying a key-value pair

Users can send a HTTP GET request to any of leader and follower nodes to query a value corresponding to a key. 
```bash
curl -X GET localhost:9091/v0/kv/key/aGVsbG8=
```


## Development

### Building the project

All data structures for gRPC and REST API are defined using [Protocol Buffers](https://developers.google.com/protocol-buffers) in the [`proto/dkv`](./proto/dkv) directory.
After adding/editing `*.proto` files in that directory, you can run the following commands to generate `*.pb.go` and `*.pb.gw.go` files in the [`pb/dkv`](./pb/dkv) directory.
```bash
make proto-clean proto-gen
```

If you want to build a `dkv` binary,
```bash
make build
```
This will build a `dkv` binary in the `./build` directory.

### Development Roadmap

| Topic                                                     | Status |
|-----------------------------------------------------------|--------|
| Basic storage features using BoltDB                       | Done   |
| Basic (Non-fault-tolerant) replication with in-memory WAL | WIP    |
| Persistent WAL                                            | TBD    |
| Fault-tolerant replication (follower failures)            | TBD    | 
| Fault-tolerant replication (leader failures)              | TBD    | 
| Dynamic cluster configuration (adding/removing followers) | TBD    | 
| WAL compaction                                            | TBD    |


## License

[Apache License 2.0](./LICENSE)