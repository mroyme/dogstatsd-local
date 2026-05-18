# Protocol Reference

dogstatsd-local implements the [DogStatsD protocol](https://docs.datadoghq.com/developers/dogstatsd/) from Datadog, which extends the [statsd protocol](https://github.com/b/statsd_spec).

## Metrics

```
metric.name:value|type|@sample_rate|#tag1:value,tag2
```

### Supported Types

| Type | Symbol | Description |
|------|--------|-------------|
| Count | `c` | Counter, tracks cumulative value |
| Gauge | `g` | Point-in-time value |
| Set | `s` | Counts unique occurrences |
| Timer | `ms` | Timing in milliseconds |
| Histogram | `h` | Statistical distribution |
| Distribution | `d` | Distribution across hosts |

### Examples

```
# Counter with tags and sample rate
page.views:1|c|@0.5|#env:dev,country:us

# Gauge
fuel.level:0.5|g

# Timer
request.time:150|ms

# Histogram
song.length:240|h|@0.5

# Distribution
page.views:42|d|#env:dev

# Set
users.uniques:1234|s
```

### Namespace Splitting

If a metric name contains a dot, the first segment becomes the namespace and the rest becomes the name:

```
page.views:1|c    → namespace="page", name="views"
fuel:0.5|g        → namespace="", name="fuel"
```

## Service Checks

```
_sc|name|status|#tags|d:timestamp|h:hostname|m:message
```

### Status Values

| Value | Meaning |
|-------|---------|
| `0` | OK |
| `1` | WARN |
| `2` | CRIT |
| `3` | UNKNOWN |

### The `m:` field

The `m:` field (message) **must** be positioned last in the datagram because its value can contain pipe (`|`) characters. The parser extracts it before splitting on `|`.

### Examples

```
# Critical with tags and hostname
_sc|Redis connection|2|#env:dev|h:db1.example.com

# Warning with a message containing pipes
_sc|db_check|1|#env:prod|m:Error: timeout|retrying

# OK with timestamp
_sc|cache_check|0|#env:staging|d:1656581400|h:cache1|m:Healthy
```

## Events

```
_e{<TITLE_LEN>,<TEXT_LEN>}:<TITLE>|<TEXT>|d:<TS>|h:<HOST>|p:<PRIORITY>|t:<ALERT_TYPE>|k:<AGG_KEY>|s:<SRC_TYPE>|#<TAGS>
```

Title and text are extracted using byte-length matching, so they can contain pipe characters and other special characters.

### Event Fields

| Field | Required | Description |
|-------|----------|-------------|
| Title | Yes | Event title |
| Text | Yes | Event body text |
| `d:` | No | Unix timestamp |
| `h:` | No | Hostname |
| `p:` | No | Priority: `normal` (default) or `low` |
| `t:` | No | Alert type: `info` (default), `warning`, `error`, `success` |
| `k:` | No | Aggregation key |
| `s:` | No | Source type |
| `#` | No | Tags (comma-separated) |

### Examples

```
# Simple event
_e{5,4}:title|text

# Event with all fields
_e{5,17}:title|Cannot parse JSON|h:host1|p:low|t:error|k:aggkey1|s:source1|#env:prod,region:us

# Event with pipe characters in text
_e{6,15}:title1|text with pipes|t:warning|#err_type:bad_file
```

## Multi-Message Datagrams

Per the statsd/DogStatsD protocol, multiple messages can be sent in a single UDP datagram, separated by newlines:

```bash
printf "page.views:1|c\nfuel.level:0.5|g" | nc -u -w1 localhost 8125
```

Both messages will be parsed and displayed separately. Carriage return (`\r`) characters are automatically stripped.