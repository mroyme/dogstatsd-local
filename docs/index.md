# dogstatsd-local

A local DogStatsD protocol inspector for debugging metrics, service checks, and events.

## Why?

Datadog is great for production metric aggregation. **dogstatsd-local** lets you inspect and debug metrics _before_ sending them to Datadog — no account required, no agent needed, no metrics polluted.

## Quick Start

Install with Go:

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

## Output Formats

### Pretty (default)

![Pretty format output](assets/pretty-format.png)

### JSON

Machine-readable, one JSON object per line:

```json
{"namespace":"page","name":"views","path":"page.views","value":1,"sample_rate":1,"tags":["env:dev"]}
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

## Features

- **All DogStatsD message types** — metrics, service checks, events
- **4 output formats** — pretty, json, short, raw
- **Forwarding** — proxy datagrams to an upstream DogStatsD server
- **Catppuccin colors** — pretty format adapts to light/dark terminal themes
- **Zero dependencies** — single binary, no Datadog account needed

## Acknowledgments

Started as a fork of [jonmorehouse/dogstatsd-local](https://github.com/jonmorehouse/dogstatsd-local), which was no longer receiving updates. Since then, this project has diverged significantly — adding service check and event support, multiple output formats, Catppuccin-themed colors, metric forwarding, and more.