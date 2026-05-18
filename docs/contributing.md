# Contributing

## Prerequisites

Install [mise](https://mise.jdx.dev/) — it will manage Go and uv automatically:

## Tasks

All common tasks are defined in `mise.toml` and run via `mise run <task>`:

| Task                   | Description                                    |
| ---------------------- | ---------------------------------------------- |
| `mise run build`       | Build the binary                               |
| `mise run test`        | Run all tests (with race detector)             |
| `mise run lint`        | Run golangci-lint                              |
| `mise run fmt`         | Format Go code (gofumpt)                       |
| `mise run fix`         | Run `go fix` on all packages                   |
| `mise run vet`         | Run `go vet`                                   |
| `mise run check`       | Run all checks (fmt → fix → vet → lint → test) |
| `mise run docs`        | Build docs site                                |
| `mise run docs-serve`  | Serve docs locally with live reload            |
| `mise run docs-deploy` | Deploy docs to GitHub Pages                    |
| `mise run clean`       | Remove `bin/` and `site/`                      |

## Development Workflow

Make a change, then run all checks:

```bash
mise run check
```

Or run individual steps:

```bash
mise run fmt
mise run test
mise run lint
```

## Documentation

Documentation lives in `docs/` and is built with [MkDocs Material](https://squidfunk.github.io/mkdocs-material/).

To preview changes locally:

```bash
mise run docs-serve
```

This starts a live-reloading server at http://127.0.0.1:8000.

Python dependencies are declared in `pyproject.toml` and managed by uv — no manual pip installs needed.
