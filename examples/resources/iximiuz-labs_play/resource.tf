resource "iximiuz-labs_play" "example" {
  playground = iximiuz-labs_playground.example.name

  # Required for community-authored playgrounds
  # safety_disclaimer_consent = true

  # Override init condition values
  # init_conditions = {
  #   runtime = "docker"
  # }
}

output "play_url" {
  value = iximiuz-labs_play.example.page_url
}
