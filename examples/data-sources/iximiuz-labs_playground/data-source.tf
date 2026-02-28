data "iximiuz-labs_playground" "docker" {
  name = "docker"
}

output "playground_title" {
  value = data.iximiuz-labs_playground.docker.title
}
