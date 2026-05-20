# dogstatsd-local

[![CI](https://github.com/mroyme/dogstatsd-local/actions/workflows/ci.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/ci.yml)
[![Trivy](https://github.com/mroyme/dogstatsd-local/actions/workflows/trivy.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/trivy.yml)
[![Release](https://github.com/mroyme/dogstatsd-local/actions/workflows/release.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/release.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/mroyme/dogstatsd-local?logo=docker)](https://hub.docker.com/r/mroyme/dogstatsd-local)

> A local DogStatsD protocol inspector for debugging metrics, service checks, and events

> [!TIP]
> Full documentation available [here](https://mroyme.github.io/dogstatsd-local/)

## Features

- **All DogStatsD message types** — metrics (count, gauge, set, timer, histogram, distribution), service checks, events
- **Multiple output formats** — pretty, json, short, raw
- **Metric forwarding** — proxy datagrams to an upstream DogStatsD server while inspecting locally
- **Catppuccin colors** — pretty format adapts to your terminal's light/dark theme

## Quick Start

Install with Homebrew:

```bash
brew install mroyme/tap/dogstatsd-local
```

Or with Go:

```bash
go install github.com/mroyme/dogstatsd-local/cmd/dogstatsd-local@latest
```

Or run with Docker:

```bash
docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local
```

Or download a [prebuilt binary](https://github.com/mroyme/dogstatsd-local/releases/latest) for Linux, macOS, or Windows (x86-64 and ARM64).

Then start it up and point your service at it:

```bash
export DD_AGENT_HOST=127.0.0.1
export DD_DOGSTATSD_PORT=8125
```

Or test manually with netcat:

```bash
# Metric
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125

# Service check
printf "_sc|Redis connection|2|#env:dev|m:Timeout" | nc -u -w1 localhost 8125

# Event
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

### Pretty (default)

![Pretty format output](https://mroyme.github.io/dogstatsd-local/assets/pretty-format.png)

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

> Application → dogstatsd-local (8126) → Datadog agent (8125)

```bash
dogstatsd-local -port 8126 -forward 127.0.0.1:8125
```

All datagrams are forwarded as-is — no modification, no data loss.

## Documentation

Full documentation is available at [mroyme.github.io/dogstatsd-local](https://mroyme.github.io/dogstatsd-local/).

- [Installation](https://mroyme.github.io/dogstatsd-local/installation/)
- [Configuration](https://mroyme.github.io/dogstatsd-local/configuration/)
- [Connecting Your Service](https://mroyme.github.io/dogstatsd-local/connecting/)
- [Output Formats](https://mroyme.github.io/dogstatsd-local/output-formats/)
- [Forwarding](https://mroyme.github.io/dogstatsd-local/forwarding/)
- [Protocol Reference](https://mroyme.github.io/dogstatsd-local/protocol/)
- [Contributing](https://mroyme.github.io/dogstatsd-local/contributing/)

## Acknowledgments

Started as a fork of [jonmorehouse/dogstatsd-local](https://github.com/jonmorehouse/dogstatsd-local), which was no longer receiving updates. Since then, this project has diverged significantly — adding service check and event support, multiple output formats, Catppuccin-themed colors, metric forwarding, and more.
