# dogstatsd-local

[![CI](https://github.com/mroyme/dogstatsd-local/actions/workflows/ci.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/ci.yml)
[![Docker](https://github.com/mroyme/dogstatsd-local/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/docker-publish.yml)
[![Trivy](https://github.com/mroyme/dogstatsd-local/actions/workflows/trivy.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/trivy.yml)
[![SLSA Go releaser](https://github.com/mroyme/dogstatsd-local/actions/workflows/go-ossf-slsa3-publish.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/go-ossf-slsa3-publish.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/mroyme/dogstatsd-local?logo=docker)](https://hub.docker.com/r/mroyme/dogstatsd-local)

> A local DogStatsD protocol inspector for debugging metrics, service checks, and events

```
dogstatsd-local
```

![Pretty format output](https://raw.githubusercontent.com/mroyme/dogstatsd-local/main/docs/assets/pretty-format.png)

## Features

- 📊 **All DogStatsD message types** — metrics (count, gauge, set, timer, histogram, distribution), service checks, events
- 🎨 **4 output formats** — pretty, json, short, raw
- 🔀 **Metric forwarding** — proxy datagrams to an upstream DogStatsD server while inspecting locally
- 🌈 **Catppuccin colors** — pretty format adapts to your terminal's light/dark theme

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

Download from the [releases page](https://github.com/mroyme/dogstatsd-local/releases/latest) for Linux, macOS, and Windows (x86-64 and ARM64).

## Quick Start

```bash
# Start listening on port 8125 (default)
dogstatsd-local

# Send a metric
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125

# Send a service check
printf "_sc|Redis connection|2|#env:dev|m:Timeout" | nc -u -w1 localhost 8125

# Send an event
printf "_e{21,21}:An exception occurred|Cannot parse CSV file|t:warning|#err_type:bad_file" | nc -u -w1 localhost 8125
```

## Flags

| Flag               | Default      | Description                                       |
| ------------------ | ------------ | ------------------------------------------------- |
| `-host`            | `0.0.0.0`    | Bind address                                      |
| `-port`            | `8125`       | UDP listen port                                   |
| `-out`             | `pretty`     | Output format: `pretty`, `json`, `short`, `raw`   |
| `-forward`         | _(disabled)_ | Forward datagrams to an upstream DogStatsD server |
| `-tags`            | _(empty)_    | Extra tags to append, comma-delimited             |
| `-max-name-width`  | `50`         | Max name length for pretty format                 |
| `-max-value-width` | `15`         | Max value length for pretty format                |
| `-debug`           | `false`      | Enable debug logging                              |

## Output Formats

### Pretty

Colorized, human-readable with Catppuccin-themed colors:

```
COUNT      page | views                                      1.00           env:dev
CRIT       Redis connection                                  Timeout env:dev
WARN       An exception occurred                             Cannot parse CSV file err_type:bad_file
```

### JSON

Machine-readable, one JSON object per line. Pipe through `jq` for pretty printing:

```json
{
  "namespace": "page",
  "name": "views",
  "path": "page.views",
  "value": 1,
  "sample_rate": 1,
  "tags": ["env:dev"]
}
```

### Short

Compact human-readable:

```
metric:count|page.views|1.00 env:dev
service_check:Redis connection|CRIT|msg:Timeout env:dev
event:title|text|priority:normal|alert:info
```

### Raw

Passthrough of the original datagram:

```
page.views:1|c|#env:dev
```

## Forwarding

Run as a middleware between your app and the Datadog agent:

```
Application → dogstatsd-local (8126) → Datadog agent (8125)
```

```bash
dogstatsd-local -port 8126 -forward 127.0.0.1:8125
```

All datagrams are forwarded as-is — no modification, no data loss.

## Documentation

Full documentation is available at [mroyme.github.io/dogstatsd-local](https://mroyme.github.io/dogstatsd-local/).

- [Installation](https://mroyme.github.io/dogstatsd-local/installation/)
- [Configuration](https://mroyme.github.io/dogstatsd-local/configuration/)
- [Output Formats](https://mroyme.github.io/dogstatsd-local/output-formats/)
- [Forwarding](https://mroyme.github.io/dogstatsd-local/forwarding/)
- [Protocol Reference](https://mroyme.github.io/dogstatsd-local/protocol/)
- [Contributing](https://mroyme.github.io/dogstatsd-local/contributing/)

## Acknowledgments

Started as a fork of [jonmorehouse/dogstatsd-local](https://github.com/jonmorehouse/dogstatsd-local), which was no longer receiving updates. Since then, this project has diverged significantly — adding service check and event support, multiple output formats, Catppuccin-themed colors, metric forwarding, and more.
