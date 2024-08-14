# Dogstatsd Local

[![Docker](https://github.com/mroyme/dogstatsd-local/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/docker-publish.yml)
[![SLSA Go releaser](https://github.com/mroyme/dogstatsd-local/actions/workflows/go-ossf-slsa3-publish.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/go-ossf-slsa3-publish.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/mroyme/dogstatsd-local?logo=docker)](https://hub.docker.com/r/mroyme/dogstatsd-local)

> A local implementation of the dogstatsd protocol from [Datadog](https://www.datadog.com)


## Why?

[Datadog](https://www.datadog.com) is great for production application metric aggregation. This project was inspired by the need to inspect and debug metrics _before_ sending them to `datadog`.

`dogstatsd-local` is a small program which understands the `dogstatsd` and `statsd` protocols. It listens on a local UDP server and writes metrics, events and service checks per the [dogstatsd protocol](https://docs.datadoghq.com/guides/dogstatsd/) to `stdout` in user configurable formats.

This can be helpful for _debugging_ metrics themselves, and to prevent polluting datadog with noisy metrics from a development environment. **dogstatsd-local** can also be used to pipe metrics as json to other processes for further processing.

## Usage

### Install with Go

```
go install github.com/mroyme/dogstatsd-local/cmd/dogstatsd-local@latest
```

### Build Manually

This is a go application with no external dependencies. Building should be as simple as running `go build` in the source directory.

Once compiled, the `dogstatsd-local` binary can be run directly:
```bash
$ ./dogstatsd-local -port=8126
```

### Prebuilt Binaries

Pre-built binaries for Linux, Mac and Windows are available for x86-64 and AArch64.
Check out the [releases](https://github.com/mroyme/dogstatsd-local/releases/latest) page.


### Docker

```bash
$ docker run -t -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local
```

## Sample Formats

### Pretty 

'Pretty' is the default format. When writing a metric such as:

```bash
$ printf "namespace.metric:1|c|#test" | nc -cu  localhost 8125
```

Running **dogstatsd-local** with the `-format raw` flag will output the plain udp packet:

```bash
$ docker run -t -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local -format pretty
COUNTER    namespace | metric                                1.00           TAGS = test
```

The output will be colored if your shell supports colors.
If colors aren't displayed properly, ensure that `TERM` is set correctly in your environment.

Pretty supports the following extra flags:
- `-max-name-width` (integer): Maximum length of name. Change if name is truncated (default 50)
- `-max-name-width` (integer): Maximum length of name. Change if value is truncated (default 50)


### Raw (no formatting)

When writing a metric such as:

```bash
$ printf "namespace.metric:1|c|#test" | nc -cu  localhost 8125
```

Running **dogstatsd-local** with the `-format raw` flag will output the plain udp packet:

```bash
$ docker run -p 8125:8125/udp mroyme/dogstatsd-local -format raw
2017/12/03 23:11:31 namespace.metric.name:1|c|@1.00|#tag1
```

### Short 

When writing a metric such as:

```bash
$ printf "namespace.metric:1|c|#test" | nc -cu  localhost 8125
```

Running **dogstatsd-local** with the `-format short` flag will output a short, albeit still human-readable metric:

```bash
$ docker run -p 8125:8125/udp mroyme/dogstatsd-local -format short
metric:counter|namespace.metric|1.00  test

```

### JSON

When writing a metric such as:
```bash
$ printf "namespace.metric:1|c|#test|extra" | nc -cu  localhost 8125
```

Running **dogstatsd-local** with the `-format json` flag will output json:

```bash
$ docker run -p 8125:8125/udp mroyme/dogstatsd-local -format json | jq .
{"namespace":"namespace","name":"metric","path":"namespace.metric","value":1,"extras":["extra"],"sample_rate":1,"tags":["test"]}
```

**dogstatsd-local** can be piped to any process that understands json via stdin. For example, to pretty print JSON with [jq](https://stedolan.github.io/jq/):

```bash
$ docker run -p 8125:8125/udp mroyme/dogstatsd-local -format json | jq .
{
  "namespace": "namespace",
  "name": "metric",
  "path": "namespace.metric",
  "value": 1,
  "extras": [
    "extra"
  ],
  "sample_rate": 1,
  "tags": [
    "test"
  ]
}
```

## TODO

- [ ] support datadog service checks
- [ ] support datadog events
- [ ] support interval aggregation of percentiles
