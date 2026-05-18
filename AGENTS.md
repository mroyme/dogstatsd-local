# Agent Development Guide

A file for [guiding coding agents](https://agents.md/).

## Commands

- **Build:** `go build -o bin/dogstatsd-local ./cmd/dogstatsd-local/main.go`
- **Run:** `./bin/dogstatsd-local -port=8125 -out=pretty`
- **Test:** `go test -v ./...`
- **Vet:** `go vet ./...`
- **Docker:** `docker run -it -e "TERM=$TERM" -p 8125:8125/udp mroyme/dogstatsd-local`

## Project Structure

- Entry point: `cmd/dogstatsd-local/main.go`
- Protocol parsing + types: `internal/messages/`
- UDP server: `internal/server/`
- Output formats: `internal/format/` (pretty, json, raw, short)

## Code Conventions

- Go 1.26+, standard `testing` package only
- Parse errors in format handlers are logged but return `nil` (never kill the pool worker)
- `m:` field in service checks must be extracted before splitting on `|` (its value can contain pipes)

## DogStatsD Protocol

- Metrics: `name:value|type|@sample_rate|#tags`
- Service checks: `_sc|name|status|#tags|d:timestamp|h:hostname|m:message` (`m:` must be last)
- Events: `_e{...}` — not yet implemented

## TODO

- [ ] Datadog event support (`_e{` prefix)
- [ ] Interval aggregation of percentiles
- [ ] Add meaningful server tests
