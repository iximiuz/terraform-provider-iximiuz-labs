data "iximiuz-labs_playgrounds" "all" {
  filter = "all"
}

output "playground_names" {
  value = [for pg in data.iximiuz-labs_playgrounds.all.playgrounds : pg.name]
}
