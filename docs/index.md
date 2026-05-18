# dogstatsd-local

A local DogStatsD protocol inspector. Listen on a UDP socket, parse metrics, service checks, and events, and print them to stdout in your choice of format.

## Why?

Datadog is great for production metric aggregation. **dogstatsd-local** lets you inspect and debug metrics _before_ sending them to Datadog — no account required, no agent needed, no metrics polluted.

## Quick Start

```bash
go install github.com/mroyme/dogstatsd-local/cmd/dogstatsd-local@latest
dogstatsd-local
```

```bash
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125
```

```
COUNT      page | views                                      1.00           env:dev
```

## Features

- **All DogStatsD message types** — metrics, service checks, events
- **4 output formats** — pretty, json, short, raw
- **Forwarding** — proxy datagrams to an upstream DogStatsD server
- **Catppuccin colors** — pretty format adapts to light/dark terminal themes
- **Zero dependencies** — single binary, no Datadog account needed

## Acknowledgments

Started as a fork of [jonmorehouse/dogstatsd-local](https://github.com/jonmorehouse/dogstatsd-local), which was no longer receiving updates. Since then, this project has diverged significantly — adding service check and event support, multiple output formats, Catppuccin-themed colors, metric forwarding, and more.