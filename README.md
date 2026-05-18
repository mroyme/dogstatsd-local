# Dogstatsd Local

[![Docker](https://github.com/mroyme/dogstatsd-local/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/docker-publish.yml)
[![Trivy](https://github.com/mroyme/dogstatsd-local/actions/workflows/trivy.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/trivy.yml)
[![SLSA Go releaser](https://github.com/mroyme/dogstatsd-local/actions/workflows/go-ossf-slsa3-publish.yml/badge.svg)](https://github.com/mroyme/dogstatsd-local/actions/workflows/go-ossf-slsa3-publish.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/mroyme/dogstatsd-local?logo=docker)](https://hub.docker.com/r/mroyme/dogstatsd-local)

> A local implementation of the dogstatsd protocol from [Datadog](https://www.datadog.com)
>
> Up-to-date fork of [jonmorehouse/dogstatsd-local](https://github.com/jonmorehouse/dogstatsd-local)


## Why?

[Datadog](https://www.datadog.com) is great for production application metric aggregation. This project was inspired by the need to inspect and debug metrics _before_ sending them to `datadog`.

`dogstatsd-local` is a small program which understands the `dogstatsd` and `statsd` protocols. It listens on a local UDP server and writes metrics, events and service checks per the [dogstatsd protocol](https://docs.datadoghq.com/guides/dogstatsd/) to `stdout` in user configurable formats.

This can be helpful for _debugging_ metrics themselves, and to prevent polluting datadog with noisy metrics from a development environment. **dogstatsd-local** can also be used to pipe metrics as json to other processes for further processing.

## Usage

### Install with Go

```
$ go install github.com/mroyme/dogstatsd-local/cmd/dogstatsd-local@latest
```

### Build Manually

Run the following command in the source directory.
```bash
$ go build -o bin/dogstatsd-local ./cmd/dogstatsd-local/main.go
```

Once compiled, the `dogstatsd-local` binary can be run directly:
```bash
$ ./bin/dogstatsd-local -port=8126
```

### Prebuilt Binaries

Pre-built binaries for Linux, Mac and Windows are available for x86-64 and AArch64.
Check out the [releases](https://github.com/mroyme/dogstatsd-local/releases/latest) page.


### Docker

```bash
$ docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local
```

## Sample Formats

### Pretty 

'Pretty' is the default format. When writing a metric such as:

```bash
$ printf "namespace.metric:1|c|#test" | nc -cu  localhost 8125
```

Running **dogstatsd-local** with the `-out pretty` flag will parse the UDP packets and print a colorized output:

```bash
$ docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local -out pretty
COUNT        namespace | metric                                1.00            test
```

When sending a service check:

```bash
$ printf "_sc|Redis connection|2|#env:dev|m:Redis connection timed out after 10s" | nc -cu  localhost 8125
```

The output is colorized by status (green=OK, yellow=WARN, red=CRIT):

```bash
CRIT        Redis connection                       Redis connection timed out after 10s  env:dev
```

When sending an event:

```bash
$ printf "_e{21,36}:An exception occurred|Cannot parse CSV file from 10.0.0.17|t:warning|#err_type:bad_file" | nc -cu  localhost 8125
```

The output is colorized by alert type (blue=info, yellow=warning, red=error, green=success):

```bash
WARN        An exception occurred                  Cannot parse CSV file from 10.0.0.17  err_type:bad_file
```

The output will be colored if your shell supports colors.
If colors aren't displayed properly, ensure that `TERM` is set correctly in your environment.

Pretty supports the following extra flags:
- `-max-name-width` (integer): Maximum length of name. Change if name is truncated (default 50)
- `-max-value-width` (integer): Maximum length of value. Change if value is truncated (default 50)
- `-debug` (boolean): Enable debug mode (default `false`)


### Raw (no formatting)

When writing a metric such as:

```bash
$ printf "namespace.metric:1|c|#test" | nc -cu  localhost 8125
```

Running **dogstatsd-local** with the `-out raw` flag will output the plain udp packet:

```bash
$ docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local -out raw
2017/12/03 23:11:31 namespace.metric.name:1|c|@1.00|#tag1
```

### Short 

When writing a metric such as:

```bash
$ printf "namespace.metric:1|c|#test" | nc -cu  localhost 8125
```

Running **dogstatsd-local** with the `-out short` flag will output a short, albeit still human-readable metric:

```bash
$ docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local -out short
metric:counter|namespace.metric|1.00  test
```

When sending a service check:

```bash
$ printf "_sc|Redis connection|2|#env:dev|m:Redis connection timed out after 10s" | nc -cu  localhost 8125
```

```bash
service_check:Redis connection|CRIT|msg:Redis connection timed out after 10s env:dev
```

When sending an event:

```bash
$ printf "_e{21,36}:An exception occurred|Cannot parse CSV file from 10.0.0.17|t:warning|#err_type:bad_file" | nc -cu  localhost 8125
```

```bash
event:An exception occurred|Cannot parse CSV file from 10.0.0.17|priority:normal|alert:warning err_type:bad_file
```

### JSON

When writing a metric such as:
```bash
$ printf "namespace.metric:1|c|#test|extra" | nc -cu  localhost 8125
```

Running **dogstatsd-local** with the `-out json` flag will output json:

```bash
$ docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local -out json | jq .
{"namespace":"namespace","name":"metric","path":"namespace.metric","value":1,"extras":["extra"],"sample_rate":1,"tags":["test"]}
```

When sending a service check:

```bash
$ printf "_sc|Redis connection|2|#env:dev|m:Redis connection timed out after 10s" | nc -cu  localhost 8125
```

```bash
$ docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local -out json | jq .
{
  "name": "Redis connection",
  "status": "CRIT",
  "message": "Redis connection timed out after 10s",
  "tags": [
    "env:dev"
  ],
  "timestamp": 1656581400
}
```

When sending an event:

```bash
$ printf "_e{21,36}:An exception occurred|Cannot parse CSV file from 10.0.0.17|t:warning|#err_type:bad_file" | nc -cu  localhost 8125
```

```bash
$ docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local -out json | jq .
{
  "title": "An exception occurred",
  "text": "Cannot parse CSV file from 10.0.0.17",
  "priority": "normal",
  "alert_type": "warning",
  "tags": [
    "err_type:bad_file"
  ],
  "timestamp": 1656581400
}
```

**dogstatsd-local** can be piped to any process that understands json via stdin. For example, to pretty print JSON with [jq](https://stedolan.github.io/jq/):

```bash
$ docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local -out json | jq .
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
