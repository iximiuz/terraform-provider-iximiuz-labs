data "iximiuz-labs_plays" "all" {}

output "running_plays" {
  value = [for p in data.iximiuz-labs_plays.all.plays : p.id]
}
