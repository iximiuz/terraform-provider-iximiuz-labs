data "iximiuz-labs_me" "current" {}

output "user_id" {
  value = data.iximiuz-labs_me.current.id
}

output "has_premium" {
  value = data.iximiuz-labs_me.current.has_premium
}
