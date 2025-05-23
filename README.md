# Reddio
[Reddio](https://www.reddio.com/) is a high performance parallel Ethereum-compatible Layer 2, leveraging
zero-knowledge technology to achieve unrivaled computation scale with
Ethereum-level security.

## Build & Run

### Prerequisites

- go 1.23.0

### Ethereum compatible

- go-ethereum v1.14.0

### Source code Build & Run

```shell
git clone git@github.com:reddio-com/reddio.git
cd reddio && make build

./reddio
```

### Docker Pull & Run

```shell
docker pull ghcr.io/reddio-com/reddio:latest
docker-compose up
```

### Check pprof

http://localhost:10199/debug/pprof