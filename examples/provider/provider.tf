# Copyright iximiuz Labs 2026
# SPDX-License-Identifier: MPL-2.0

# Configure the iximiuz Labs provider.
#
# Credentials can be supplied here, via environment variables
# (IXIMIUZ_SESSION_ID, IXIMIUZ_ACCESS_TOKEN), or by installing
# labctl and running `labctl auth login` (reads ~/.iximiuz/labctl/config.yaml).
provider "iximiuz-labs" {
  # api_url      = "https://labs.iximiuz.com/api"  # optional, defaults shown
  # session_id   = "..."                            # or IXIMIUZ_SESSION_ID
  # access_token = "..."                            # or IXIMIUZ_ACCESS_TOKEN
}
