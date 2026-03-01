# Terraform Provider for iximiuz Labs

The `iximiuz-labs` Terraform provider manages resources on the
[iximiuz Labs](https://labs.iximiuz.com) platform.

## Installation

```hcl
terraform {
  required_providers {
    iximiuz-labs = {
      source = "iximiuz/iximiuz-labs"
    }
  }
}
```

## Resources

| Resource | Description |
|---|---|
| `iximiuz-labs_playground` | Create and manage playground definitions (machines, networks, tabs, etc.) |
| `iximiuz-labs_play` | Launch a play session from a playground |
| `iximiuz-labs_play_port` | Forward a port in a running play |
| `iximiuz-labs_play_shell` | Open a shell session in a running play |

## Data Sources

| Data Source | Description |
|---|---|
| `iximiuz-labs_playground` | Read a single playground by name |
| `iximiuz-labs_playgrounds` | List all playgrounds |
| `iximiuz-labs_play` | Read a single play by ID |
| `iximiuz-labs_plays` | List plays, optionally filtered by playground |
| `iximiuz-labs_me` | Read the currently authenticated user |

## Authentication

Credentials are resolved in order:

1. Provider configuration (`session_id` / `access_token` attributes)
2. Environment variables (`IXIMIUZ_SESSION_ID` / `IXIMIUZ_ACCESS_TOKEN`)
3. labctl config file (`~/.iximiuz/labctl/config.yaml`)

The easiest approach is to install [labctl](https://github.com/iximiuz/labctl)
and run `labctl auth login` — the provider will pick up the saved credentials
automatically.

## Example Usage

```hcl
provider "iximiuz-labs" {}

# Look up an existing playground
data "iximiuz-labs_playground" "docker" {
  name = "docker"
}

# Launch a play session
resource "iximiuz-labs_play" "session" {
  playground = data.iximiuz-labs_playground.docker.name
}

# Retrieve current user info
data "iximiuz-labs_me" "current" {}
```

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25 (only for building from source)

## Building from Source

```shell
git clone https://github.com/iximiuz/terraform-provider-iximiuz-labs.git
cd terraform-provider-iximiuz-labs
go install
```

## Developing

```shell
# Run acceptance tests (requires valid API credentials)
make testacc

# Regenerate documentation
make generate
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, running tests, and
submitting pull requests.

## License

MPL-2.0 — Copyright iximiuz Labs 2026
