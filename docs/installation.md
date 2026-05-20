# Installation

## Homebrew

```bash
brew install mroyme/tap/dogstatsd-local
```

## Go

```bash
go install github.com/mroyme/dogstatsd-local/cmd/dogstatsd-local@latest
```

## Docker

### Docker Hub

```bash
docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local
```

### GHCR

```bash
docker run -it -e "TERM=$TERM" -p 8125:8125/udp ghcr.io/mroyme/dogstatsd-local
```

Pass CLI flags via the command:

```bash
docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local -out json -tags env:dev
```

## Prebuilt Binaries

Download from the [releases page](https://github.com/mroyme/dogstatsd-local/releases/latest) for Linux, macOS, and Windows (x86-64 and ARM64).

## Build from Source

```bash
go build -o bin/dogstatsd-local ./cmd/dogstatsd-local/main.go
```