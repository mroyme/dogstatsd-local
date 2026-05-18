# Configuration

All options are CLI flags. Run `dogstatsd-local -h` to see defaults.

## Flags

| Flag               | Default      | Description                                                                   |
| ------------------ | ------------ | ----------------------------------------------------------------------------- |
| `-host`            | `0.0.0.0`    | Bind address                                                                  |
| `-port`            | `8125`       | UDP listen port                                                               |
| `-out`             | `pretty`     | Output format: `pretty`, `json`, `short`, `raw`                               |
| `-forward`         | _(disabled)_ | Forward raw datagrams to an upstream DogStatsD server (e.g. `127.0.0.1:8126`) |
| `-tags`            | _(empty)_    | Extra tags to append, comma-delimited                                         |
| `-max-name-width`  | `50`         | Max name length for pretty format (values below 50 have no effect)            |
| `-max-value-width` | `15`         | Max value length for pretty format                                            |
| `-debug`           | `false`      | Enable debug logging                                                          |

## Examples

Default pretty output on port 8125:

```shell
dogstatsd-local
```

JSON output with extra tags:

```shell
dogstatsd-local -out json -tags env:dev,team:platform
```

Listen on a custom port and forward to Datadog agent:

```bash
dogstatsd-local -port 8126 -forward 127.0.0.1:8125
```

Custom column widths for long metric names:

```bash
dogstatsd-local -max-name-width 80 -max-value-width 20
```
