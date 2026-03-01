resource "iximiuz-labs_playground" "example" {
  name_prefix = "my-playground"
  base        = "docker"
  title       = "My Playground"
  description = "A custom playground for Docker experimentation"

  categories = ["linux"]

  tab {
    kind    = "terminal"
    name    = "Terminal"
    machine = "docker-01"
  }
}
