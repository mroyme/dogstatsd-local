# Output Formats

## Pretty (default)

![Pretty format output](assets/pretty-format.png)

Colorized, human-readable output with Catppuccin-themed colors that adapt to your terminal's light/dark mode.

```bash
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125
```

```
COUNT      page | views                                      1.00           env:dev
```

### Metrics

Metric types are color-coded:

| Type         | Label   | Color |
| ------------ | ------- | ----- |
| Counter      | `COUNT` | Green |
| Gauge        | `GAUGE` | Teal  |
| Set          | `SET`   | Pink  |
| Timer        | `TIMER` | Mauve |
| Histogram    | `HIST`  | Blue  |
| Distribution | `DISTR` | Sky   |

### Service Checks

```bash
printf "_sc|Redis connection|2|#env:dev|m:Timeout" | nc -u -w1 localhost 8125
```

```
CRIT       Redis connection                                  Timeout env:dev
```

Status is color-coded: green for OK, yellow for WARN, red for CRIT.

### Events

```bash
printf "_e{21,36}:An exception occurred|Cannot parse CSV file|t:warning|#err_type:bad_file" | nc -u -w1 localhost 8125
```

```
WARN       An exception occurred                             Cannot parse CSV file err_type:bad_file
```

Alert type is color-coded: blue for INFO, yellow for WARN, red for ERR, green for OK.

### Extra Tags

Use `-tags` to append tags to every message:

```bash
dogstatsd-local -tags env:dev,team:platform
```

### Column Widths

Use `-max-name-width` and `-max-value-width` to adjust column sizes:

```bash
dogstatsd-local -max-name-width 80 -max-value-width 20
```

---

## Short

Compact human-readable, one line per message.

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

---

## JSON

Machine-readable, one JSON object per line. Pipe through `jq` for pretty printing.

```bash
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125
```

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

```bash
printf "_sc|Redis connection|2|#env:dev|m:Timeout" | nc -u -w1 localhost 8125
```

```json
{
  "name": "Redis connection",
  "status": "CRIT",
  "message": "Timeout",
  "tags": ["env:dev"],
  "timestamp": 1656581400
}
```

```bash
printf "_e{5,4}:title|text|t:info" | nc -u -w1 localhost 8125
```

```json
{
  "title": "title",
  "text": "text",
  "priority": "normal",
  "alert_type": "info",
  "tags": []
}
```

### Piping to other tools

```bash
dogstatsd-local -out json | jq .
```

```bash
dogstatsd-local -out json | jq 'select(.type == "metric")'
```

---

## Raw

Passthrough of the original datagram — exactly what was received.

```bash
printf "page.views:1|c|#env:dev" | nc -u -w1 localhost 8125
```

```
page.views:1|c|#env:dev
```
