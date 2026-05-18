# Dogstatsd Local

[![Docker](https://github.com/mroyme/dogstatsd-local/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/docker-publish.yml)
[![Trivy](https://github.com/mroyme/dogstatsd-local/actions/workflows/trivy.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/trivy.yml)
[![SLSA Go releaser](https://github.com/mroyme/dogstatsd-local/actions/workflows/go-ossf-slsa3-publish.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/go-ossf-slsa3-publish.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/mroyme/dogstatsd-local?logo=docker)](https://hub.docker.com/r/mroyme/dogstatsd-local)

> A local implementation of the DogStatsD protocol from [Datadog](https://www.datadoghq.com)
>
> [!NOTE]
> Started as a fork of [jonmorehouse/dogstatsd-local](https://github.com/jonmorehouse/dogstatsd-local), which was no longer receiving updates. Since then, this project has diverged significantly — adding service check and event support, multiple output formats, Catppuccin-themed colors, metric forwarding, and more.

`dogstatsd-local` listens on a UDP socket, parses DogStatsD (and statsd) metric messages, and outputs them to stdout in your choice of format. Use it to inspect and debug metrics locally before sending them to Datadog.

## Table of Contents

- [Installation](#installation)
- [Flags](#flags)
- [Features](#features)
- [Output Formats](#output-formats)
  - [Pretty (default)](#pretty-default)
  - [Short](#short)
  - [JSON](#json)
  - [Raw](#raw)
- [Forwarding](#forwarding)
- [DogStatsD Protocol](#dogstatsd-protocol)

## Installation

### Go

```bash
go install github.com/mroyme/dogstatsd-local/cmd/dogstatsd-local@latest
```

### Docker

```bash
docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local
```

### Prebuilt Binaries

Download from the [releases](https://github.com/mroyme/dogstatsd-local/releases/latest) page for Linux, macOS, and Windows (x86-64 and ARM64).

### Build from Source

```bash
go build -o bin/dogstatsd-local ./cmd/dogstatsd-local/main.go
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-host` | `0.0.0.0` | Bind address |
| `-port` | `8125` | UDP listen port |
| `-out` | `pretty` | Output format: `pretty`, `json`, `short`, `raw` |
| `-forward` | (disabled) | Forward raw datagrams to an upstream DogStatsD server (e.g. `127.0.0.1:8126`) |
| `-tags` | (empty) | Extra tags to append, comma-delimited |
| `-max-name-width` | `50` | Max name length for pretty format (values below 50 have no effect) |
| `-max-value-width` | `15` | Max value length for pretty format |
| `-debug` | `false` | Enable debug logging |

## Features

- **Metrics** — counters (`c`), gauges (`g`), sets (`s`), timers (`ms`), histograms (`h`), distributions (`d`)
- **Service checks** — `_sc|name|status|#tags|h:hostname|m:message`
- **Events** — `_e{titleLen,textLen}:title|text|d:ts|h:host|p:priority|t:alert|k:aggkey|s:src_type|#tags`
- **Forwarding** — proxy all datagrams to an upstream DogStatsD server with `-forward`
- **Multi-message datagrams** — splits on `\n` per the DogStatsD protocol
- **Extra tags** — append tags to every metric with `-tags`
- **Catppuccin colors** — pretty format adapts to light/dark terminal themes

## Output Formats

### Pretty (default)

Colorized, human-readable output:

```bash
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125
```

```
COUNT      page | views                                      1.00           env:dev
```

Service checks are colorized by status (green=OK, yellow=WARN, red=CRIT):

```bash
printf "_sc|Redis connection|2|#env:dev|m:Timeout" | nc -u -w1 localhost 8125
```

```
CRIT       Redis connection                                  Timeout env:dev
```

Events are colorized by alert type (blue=INFO, yellow=WARN, red=ERR, green=OK):

```bash
printf "_e{21,36}:An exception occurred|Cannot parse CSV file|t:warning|#err_type:bad_file" | nc -u -w1 localhost 8125
```

```
WARN       An exception occurred                             Cannot parse CSV file err_type:bad_file
```

### Short

Compact human-readable:

```bash
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125
```

```
metric:count|page.views|1.00 env:dev
```

```bash
printf "_sc|Redis connection|2|#env:dev|m:Timeout" | nc -u -w1 localhost 8125
```

```
service_check:Redis connection|CRIT|msg:Timeout env:dev
```

```bash
printf "_e{5,4}:title|text|t:info" | nc -u -w1 localhost 8125
```

```
event:title|text|priority:normal|alert:info
```

### JSON

Machine-readable, one JSON object per line. Pipe through `jq` for pretty printing:

```bash
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125
```

```json
{"namespace":"page","name":"views","path":"page.views","value":1,"sample_rate":1,"tags":["env:dev"]}
```

```bash
printf "_sc|Redis connection|2|#env:dev|m:Timeout" | nc -u -w1 localhost 8125
```

```json
{"name":"Redis connection","status":"CRIT","message":"Timeout","tags":["env:dev"],"timestamp":1656581400}
```

```bash
printf "_e{5,4}:title|text|t:info" | nc -u -w1 localhost 8125
```

```json
{"title":"title","text":"text","priority":"normal","alert_type":"info","tags":[]}
```

### Raw

Passthrough of the original datagram:

```bash
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125
```

```
page.views:1|c|#env:dev
```

## Forwarding

Use `-forward` to proxy all datagrams to an upstream DogStatsD server while still inspecting them locally:

```bash
dogstatsd-local -forward 127.0.0.1:8126
```

This is useful for debugging in environments where you still want metrics to reach Datadog.

## DogStatsD Protocol

- **Metrics:** `name:value|type|@sample_rate|#tags` (types: `c` count, `g` gauge, `s` set, `ms` timer, `h` histogram, `d` distribution)
- **Service checks:** `_sc|name|status|#tags|d:timestamp|h:hostname|m:message` (`m:` must be last, can contain `|`)
- **Events:** `_e{<TITLE_LEN>,<TEXT_LEN>}:<TITLE>|<TEXT>|d:<TS>|h:<HOST>|p:<PRIORITY>|t:<ALERT_TYPE>|k:<AGG_KEY>|s:<SRC_TYPE>|#<TAGS>`
