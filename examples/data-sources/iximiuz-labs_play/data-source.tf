data "iximiuz-labs_play" "example" {
  id = "some-play-id"
}

output "play_status" {
  value = data.iximiuz-labs_play.example.status
}
