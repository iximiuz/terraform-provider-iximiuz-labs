# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```shell
make build       # Compile the provider
make install     # Build and install the provider locally
make lint        # Run golangci-lint
make fmt         # Format Go source files
make test        # Run unit tests (no credentials required)
make testacc     # Run acceptance tests (requires API credentials)
make generate    # Regenerate docs from templates
```

Run a single test package:
```shell
go test -v -run TestFunctionName ./internal/provider/
go test -v -run TestFunctionName ./internal/client/
```

Acceptance tests require environment variables:
```shell
export IXIMIUZ_SESSION_ID=...
export IXIMIUZ_ACCESS_TOKEN=...
TF_ACC=1 go test -v -timeout 120m ./internal/provider/
```

## Architecture

The provider follows the standard two-layer Terraform provider pattern:

### `internal/client/` — API client layer
- `client.go` — HTTP client with Basic Auth (`base64(session_id:access_token)`), exponential backoff retry (up to 5 retries, max 10s), and `X-Ratelimit-Reset` header handling
- `models.go` — All API request/response structs (Playground, Play, Challenge, Course, etc.)
- `errors.go` — Sentinel errors (`ErrNotFound`, `ErrAuthenticationRequired`, etc.) wrapped in `APIError`; use `client.IsNotFound(err)` to detect 404s
- One file per resource type (`playgrounds.go`, `plays.go`, `challenges.go`, etc.) containing API method implementations

### `internal/provider/` — Terraform resource/datasource layer
- `provider.go` — Credential resolution (config → env vars `IXIMIUZ_SESSION_ID`/`IXIMIUZ_ACCESS_TOKEN` → `~/.iximiuz/labctl/config.yaml`); registers all resources and data sources
- `poll.go` — `waitForPlayState()` helper that polls the play API every 2s until a target state is reached
- One `<name>_resource.go` and `<name>_data_source.go` file per resource/datasource, each with a corresponding `_test.go`

### Key patterns
- Resources use Terraform Plugin Framework (not SDKv2). Schema attributes use `types.String`, `types.List`, etc.
- Immutable schema fields use `planmodifier.RequiresReplace()`; computed-but-known-after-apply fields use `UseStateForUnknown()`
- The `*client.Client` is passed as `ResourceData`/`DataSourceData` in `provider.Configure()` and retrieved via type assertion in each resource's `Configure()` method
- In `Read` methods, HTTP 404 from the API should call `resp.State.RemoveResource(ctx)` to handle out-of-band deletion gracefully
- `roadmap_resource.go` is a placeholder — the API endpoint does not yet exist

### Documentation generation
Templates live in `templates/` and `examples/`. Run `make generate` (which runs `terraform-plugin-docs` from the `tools/` module) to regenerate `docs/`.
