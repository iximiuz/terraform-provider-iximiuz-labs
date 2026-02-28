# Contributing to terraform-provider-iximiuz-labs

Thank you for your interest in contributing! This document covers everything you
need to get started.

## Prerequisites

- [Go](https://golang.org/doc/install) >= 1.25
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [golangci-lint](https://golangci-lint.run/welcome/install/) (for linting)

## Setting Up the Development Environment

1. Clone the repository:

   ```shell
   git clone https://github.com/iximiuz/terraform-provider-iximiuz-labs.git
   cd terraform-provider-iximiuz-labs
   ```

2. Build the provider:

   ```shell
   go build ./...
   ```

3. Install locally:

   ```shell
   go install ./...
   ```

## Running Tests

Unit tests (no credentials required):

```shell
make test
```

Acceptance tests (requires valid iximiuz Labs API credentials):

```shell
make testacc
```

Acceptance tests hit the live API, so make sure you are authenticated via
`labctl auth login` or have `IXIMIUZ_SESSION_ID` / `IXIMIUZ_ACCESS_TOKEN`
exported in your environment.

## Running Lints

```shell
make lint
```

This runs [golangci-lint](https://golangci-lint.run/) with the configuration
defined in `.golangci.yml`.

## Regenerating Documentation

Provider documentation in `docs/` is generated from code. After changing
schemas or examples, regenerate it:

```shell
make generate
```

## Submitting Pull Requests

1. Fork the repository and create a feature branch from `main`.
2. Make your changes and ensure all tests pass (`make test`).
3. Run the linter (`make lint`) and fix any issues.
4. If you changed provider schemas or examples, run `make generate` and commit
   the updated docs.
5. Open a pull request against `main` with a clear description of the change.

## Code Style

Code style is enforced by `golangci-lint`. Run `make lint` before submitting
to catch any issues early. The project follows standard Go conventions — see
[Effective Go](https://go.dev/doc/effective_go) for guidance.

## License and Contribution Terms

By submitting a pull request, you agree that your contributions will be licensed
under the [MPL-2.0 License](LICENSE) that covers this project.

Please also review our [Code of Conduct](.github/CODE_OF_CONDUCT.md) — all
contributors are expected to follow it.
